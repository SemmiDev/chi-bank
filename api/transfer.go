package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SemmiDev/chi-bank/common"
	"github.com/SemmiDev/chi-bank/common/token"
	db "github.com/SemmiDev/chi-bank/db/sqlc"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id"`
	ToAccountID   int64  `json:"to_account_id"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
}

func (s *Server) createTransfer(w http.ResponseWriter, r *http.Request) {
	var req transferRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		Error(w, http.StatusBadRequest, ErrRequestBody)
		return
	}

	if err = validation.ValidateStruct(&req,
		validation.Field(&req.FromAccountID, validation.Required, validation.Min(1)),
		validation.Field(&req.ToAccountID, validation.Required, validation.Min(1)),
		validation.Field(&req.Amount, validation.Required, validation.Min(1)),
		validation.Field(&req.Currency, validation.Required, validation.In(common.USD, common.EUR, common.IDR, common.RM)),
	); err != nil {
		Error(w, http.StatusBadRequest, err)
		return
	}

	fromAccount, err := s.validAccount(r.Context(), req.FromAccountID, req.Currency)
	if err != nil {
		if err == sql.ErrNoRows {
			Error(w, http.StatusNotFound, err)
			return
		}
		Error(w, http.StatusForbidden, err)
		return
	}

	if fromAccount.Balance < req.Amount {
		err := errors.New("balance is not sufficient")
		Error(w, http.StatusBadRequest, err)
		return
	}

	authPayload := r.Context().Value(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("from account doesn't belong to the authenticated user")
		Error(w, http.StatusUnauthorized, err)
		return
	}

	_, err = s.validAccount(r.Context(), req.ToAccountID, req.Currency)
	if err != nil {
		if err == sql.ErrNoRows {
			Error(w, http.StatusNotFound, err)
			return
		}
		Error(w, http.StatusForbidden, err)
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := s.store.TransferTx(r.Context(), arg)
	if err != nil {
		Error(w, http.StatusInternalServerError, err)
		return
	}

	Success(w, http.StatusOK, result)
	return
}

func (s *Server) validAccount(ctx context.Context, accountID int64, currency string) (db.Account, error) {
	account, err := s.store.GetAccount(ctx, accountID)
	if err != nil {
		return account, err
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		return account, err
	}

	return account, nil
}
