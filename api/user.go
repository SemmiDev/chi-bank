package api

import (
	"encoding/json"
	"net/http"
	"time"

	db "github.com/SemmiDev/chi-bank/db/sqlc"
	"github.com/SemmiDev/chi-bank/util/password"
	"github.com/SemmiDev/chi-bank/util/validator"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
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

func (s *Server) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		MarshalError(w, http.StatusBadRequest, ErrRequestBody)
		return
	}

	err = validator.Struct(req)
	if err != nil {
		MarshalError(w, http.StatusBadRequest, err)
		return
	}

	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		MarshalError(w, http.StatusInternalServerError, err)
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
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				MarshalError(w, http.StatusForbidden, err)
				return
			}
		}
		MarshalError(w, http.StatusInternalServerError, err)
		return
	}

	rsp := newUserResponse(user)
	MarshalPayload(w, http.StatusCreated, rsp)
}
