package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/thorgnir-go-study/yp-diploma/internal/api/handler"
	"github.com/thorgnir-go-study/yp-diploma/internal/config"
	"github.com/thorgnir-go-study/yp-diploma/internal/pkg/auth"
	"github.com/thorgnir-go-study/yp-diploma/internal/usecase/order"
	"github.com/thorgnir-go-study/yp-diploma/internal/usecase/user"
	"github.com/thorgnir-go-study/yp-diploma/internal/usecase/withdrawal"
	"net/http"
)

func NewServer(cfg config.Config, userService user.UseCase, ordersService order.UseCase, withdrawalsService withdrawal.UseCase) *http.Server {
	jwtWrapper := auth.NewJwtWrapper(cfg.JWTSecret, "gophermart", 24)
	authService := handler.NewAuth(jwtWrapper)
	router := NewRouter(authService, userService, ordersService, withdrawalsService)
	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router,
	}

	return server
}

func NewRouter(authService *handler.Auth, userService user.UseCase, ordersService order.UseCase, withdrawalsService withdrawal.UseCase) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)

	// public
	r.Group(func(r chi.Router) {
		handler.MakeUserHandlers(r, userService, authService)
	})

	// protected routes
	r.Group(func(r chi.Router) {
		authService.RegisterAuthMiddleware(r)

		handler.MakeOrderHandlers(r, ordersService)
		handler.MakeWithdrawalHandlers(r, withdrawalsService, ordersService)
	})

	return r
}
