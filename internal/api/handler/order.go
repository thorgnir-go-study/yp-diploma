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
	"io"
	"net/http"
	"time"
)

type userOrder struct {
	Number     string           `json:"number"`
	Status     string           `json:"status"`
	Accrual    *decimal.Decimal `json:"accrual,omitempty"`
	UploadedAt time.Time        `json:"uploaded_at"`
}

func createOrderHandler(service order.UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := jwt.GetClaimsFromContext(r.Context())
		if err != nil || claims.UserID == entity.NilID {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		log.Info().Str("currentUser", claims.UserID.String()).Msg("CurrentUser")

		orderNumberRaw, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		orderNumber, err := entity.StringToOrderNumber(string(orderNumberRaw))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		}

		_, err = service.CreateOrder(r.Context(), orderNumber, claims.UserID)
		if err != nil {
			if errors.Is(err, entity.ErrOrderAlreadyRegistered) {
				w.WriteHeader(http.StatusOK)
				return
			}
			if errors.Is(err, entity.ErrOrderRegisteredByAnotherUser) {
				http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
				return
			}

			log.Error().Err(err).Msg("Error while creating order")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)

	}
}

func listOrdersHandler(service order.UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := jwt.GetClaimsFromContext(r.Context())
		if err != nil || claims.UserID == entity.NilID {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		//_, err = service.CreateOrder(r.Context(), orderNumber, claims.UserID)
		orders, err := service.GetUserOrders(r.Context(), claims.UserID)
		if err != nil {
			log.Error().Err(err).Msg("Error while creating order")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if len(orders) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		ordersResult := make([]userOrder, len(orders))
		for i := range orders {
			o := orders[i]
			var accrual *decimal.Decimal
			if o.Accrual.Valid {
				accrual = &o.Accrual.Decimal
			}
			ordersResult[i] = userOrder{
				Number:     o.Number.String(),
				Status:     o.Status.String(),
				Accrual:    accrual,
				UploadedAt: o.UploadedAt,
			}
		}

		respJson, err := json.Marshal(ordersResult)
		if err != nil {
			log.Error().Err(err).Msg("Error while serializing orders list")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, err = w.Write(respJson)
		if err != nil {
			log.Error().Err(err).Msg("write response failed")
		}
	}
}

func MakeOrderHandlers(r chi.Router, service order.UseCase) {
	r.Post("/api/user/orders", createOrderHandler(service))
	r.Get("/api/user/orders", listOrdersHandler(service))
}
