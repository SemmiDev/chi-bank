package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/SemmiDev/chi-bank/common"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"net/http"
	"strings"
	"time"

	db "github.com/SemmiDev/chi-bank/db/sqlc"
)

type createUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (s *Server) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		Error(w, http.StatusBadRequest, ErrRequestBody)
		return
	}

	if err = validation.ValidateStruct(&req,
		validation.Field(&req.Username, validation.Required, is.Alphanumeric),
		validation.Field(&req.Password, validation.Required, validation.Length(6,100)),
		validation.Field(&req.FullName, validation.Required),
		validation.Field(&req.Email, validation.Required, is.Email),
	); err != nil {
		Error(w, http.StatusBadRequest, err)
		return
	}

	hashedPassword, err := common.HashPassword(req.Password)
	if err != nil {
		Error(w, http.StatusInternalServerError, err)
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := s.store.CreateUser(r.Context(), arg)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			Error(w, http.StatusForbidden, errors.New("username or email already taken"))
			return
		}
		Error(w, http.StatusInternalServerError, err)
	}

	rsp := newUserResponse(user)
	Success(w, http.StatusCreated, rsp)
	return
}


type loginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (s *Server) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var req loginUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		Error(w, http.StatusBadRequest, ErrRequestBody)
		return
	}

	if err = validation.ValidateStruct(&req,
		validation.Field(&req.Username, validation.Required, is.Alphanumeric),
		validation.Field(&req.Password, validation.Required, validation.Length(6,100)),
	); err != nil {
		Error(w, http.StatusBadRequest, err)
		return
	}

	user, err := s.store.GetUser(r.Context(), req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			Error(w, http.StatusNotFound, err)
			return
		}
		Error(w, http.StatusInternalServerError, err)
		return
	}

	err = common.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		Error(w, http.StatusUnauthorized, errors.New("wrong password"))
		return
	}

	accessToken, err := s.tokenMaker.CreateToken(
		user.Username,
		s.config.AccessTokenDuration,
	)
	if err != nil {
		Error(w, http.StatusInternalServerError, err)
		return
	}

	rsp := loginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}

	Success(w, http.StatusCreated, rsp)
	return
}