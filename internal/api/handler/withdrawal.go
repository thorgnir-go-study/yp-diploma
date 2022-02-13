package handler

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"github.com/thorgnir-go-study/yp-diploma/internal/api/middleware/jwt"
	"github.com/thorgnir-go-study/yp-diploma/internal/entity"
	"github.com/thorgnir-go-study/yp-diploma/internal/usecase/order"
	"github.com/thorgnir-go-study/yp-diploma/internal/usecase/withdrawal"
	"io"
	"net/http"
	"time"
)

type userWithdrawal struct {
	OrderNumber string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

type createWithdrawalRequest struct {
	OrderNumber string          `json:"order"`
	Sum         decimal.Decimal `json:"sum"`
}

type balanceResponse struct {
	Accruals  float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

func listWithdrawalsHandler(service withdrawal.UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := jwt.GetClaimsFromContext(r.Context())
		if err != nil || claims.UserID == entity.NilID {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		withdrawals, err := service.GetWithdrawals(r.Context(), claims.UserID)

		if err != nil {
			log.Error().Err(err).Msg("Error while getting withdrawals")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if len(withdrawals) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		withdrawalsResult := make([]userWithdrawal, len(withdrawals))
		for i := range withdrawalsResult {
			wd := withdrawals[i]

			withdrawalsResult[i] = userWithdrawal{
				OrderNumber: wd.OrderNumber.String(),
				Sum:         wd.Sum.InexactFloat64(),
				ProcessedAt: wd.ProcessedAt,
			}
		}

		respJSON, err := json.Marshal(withdrawalsResult)
		if err != nil {
			log.Error().Err(err).Msg("Error while serializing withdrawals list")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, err = w.Write(respJSON)
		if err != nil {
			log.Error().Err(err).Msg("write response failed")
		}
	}
}

func createWithdrawalHandler(service withdrawal.UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := jwt.GetClaimsFromContext(r.Context())
		if err != nil || claims.UserID == entity.NilID {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		bodyContent, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			log.Error().Err(err).Msg("could not read request body")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var req createWithdrawalRequest
		if err = json.Unmarshal(bodyContent, &req); err != nil {
			log.Info().Err(err).Msg("invalid json")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		orderNumber, err := entity.StringToOrderNumber(req.OrderNumber)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		}

		err = service.CreateWithdrawal(r.Context(), claims.UserID, orderNumber, req.Sum)
		if err != nil {
			if errors.Is(err, entity.ErrInsufficientBalance) {
				http.Error(w, http.StatusText(http.StatusPaymentRequired), http.StatusPaymentRequired)
				return
			}
			log.Error().Err(err).Msg("Error while creating withdrawal")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func balanceHandler(withdrawalsService withdrawal.UseCase, ordersService order.UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := jwt.GetClaimsFromContext(r.Context())
		if err != nil || claims.UserID == entity.NilID {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		withdrawalsSum, err := withdrawalsService.GetWithdrawalsSum(r.Context(), claims.UserID)
		if err != nil {
			log.Error().Err(err).Msg("Error while getting withdrawals")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		accrualsSum, err := ordersService.GetAccrualsSum(r.Context(), claims.UserID)
		if err != nil {
			log.Error().Err(err).Msg("Error while getting accruals sum")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		result := &balanceResponse{
			Accruals:  accrualsSum.Decimal.InexactFloat64(),
			Withdrawn: withdrawalsSum.Decimal.InexactFloat64(),
		}

		respJSON, err := json.Marshal(result)
		if err != nil {
			log.Error().Err(err).Msg("Error while serializing balance response")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, err = w.Write(respJSON)
		if err != nil {
			log.Error().Err(err).Msg("write response failed")
		}
	}
}

func MakeWithdrawalHandlers(r chi.Router, withdrawalService withdrawal.UseCase, orderService order.UseCase) {
	r.Get("/api/user/withdrawals", listWithdrawalsHandler(withdrawalService))
	r.Post("/api/user/balance/withdraw", createWithdrawalHandler(withdrawalService))
	r.Get("/api/user/balance", balanceHandler(withdrawalService, orderService))
}
