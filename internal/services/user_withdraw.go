package services

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kirillshkro/gmart-loyalty/internal/types"
	"github.com/kirillshkro/gmart-loyalty/pkg/utils"
)

func (s Service) UserWithdraw(w http.ResponseWriter, r *http.Request) {
	if !s.cookieExist(r, authCookie) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	uc, err := r.Cookie(authCookie)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userID, err := s.userFromCookie(uc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ctx := context.WithValue(context.TODO(), types.UserID, userID)
	var withdraw types.WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&withdraw); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !utils.Valid(withdraw.Order) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if err = s.Repo.SetWithdraw(ctx, &withdraw); err != nil {
		if _, ok := errors.AsType[*types.ErrInsufficientBalance](err); ok {
			w.WriteHeader(http.StatusPaymentRequired)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
