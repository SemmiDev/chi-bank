package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/SemmiDev/chi-bank/common"
	"github.com/SemmiDev/chi-bank/common/token"
	db "github.com/SemmiDev/chi-bank/db/sqlc"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"
	"strconv"
	"strings"
)

type createAccountRequest struct {
	Currency string `json:"currency"`
}

func (s *Server) createAccount(w http.ResponseWriter, r *http.Request) {
	var req createAccountRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		Error(w, http.StatusBadRequest, ErrRequestBody)
		return
	}

	if err = validation.ValidateStruct(&req,
		validation.Field(&req.Currency, validation.Required, validation.In(
			common.USD,
			common.EUR,
			common.IDR,
			common.RM),
		),
	); err != nil {
		if err == validation.ErrInInvalid {
			Error(w, http.StatusBadRequest, errors.New("please input the support currency"))
			return
		}
		Error(w, http.StatusBadRequest, err)
		return
	}

	authPayload := r.Context().Value(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := s.store.CreateAccount(r.Context(), arg)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			Error(w, http.StatusForbidden, errors.New("username or email already taken"))
			return
		}
		Error(w, http.StatusInternalServerError, err)
		return
	}

	Success(w, http.StatusOK, account)
	return
}

func (s *Server) getAccount(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		Error(w, http.StatusBadRequest, err)
		return
	}

	account, err := s.store.GetAccount(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			Error(w, http.StatusNotFound, err)
			return
		}
		Error(w, http.StatusInternalServerError, err)
	}

	authPayload := r.Context().Value(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		err := errors.New("account doesn't belong to the authenticated user")
		Error(w, http.StatusUnauthorized, err)
		return
	}

	Success(w, http.StatusOK, account)
	return
}

type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (s *Server) listAccounts(w http.ResponseWriter, r *http.Request) {
	pageID, err := strconv.Atoi(r.URL.Query().Get("page_id"))
	if err != nil {
		Error(w, http.StatusBadRequest, err)
		return
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		Error(w, http.StatusBadRequest, err)
		return
	}
	if pageID == 0 || pageSize == 0 {
		Error(w, http.StatusBadRequest, errors.New("please provide page_id & page_size param"))
		return
	}
	if pageID < 1 {
		Error(w, http.StatusBadRequest, errors.New("please provide page_id greater than 1"))
		return
	}
	if pageSize < 5 || pageSize > 10 {
		Error(w, http.StatusBadRequest, errors.New("please provide page_size with range 5-10"))
		return
	}

	authPayload := r.Context().Value(authorizationPayloadKey).(*token.Payload)
	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  int32(pageSize),
		Offset: int32((pageID - 1) * pageSize),
	}

	accounts, err := s.store.ListAccounts(r.Context(), arg)
	if err != nil {
		Error(w, http.StatusInternalServerError, err)
		return
	}

	Success(w, http.StatusOK, accounts)
	return
}
