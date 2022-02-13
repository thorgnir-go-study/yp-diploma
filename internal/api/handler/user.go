package handler

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/thorgnir-go-study/yp-diploma/internal/entity"
	"github.com/thorgnir-go-study/yp-diploma/internal/usecase/user"
	"io"
	"net/http"
)

type registerRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func registerUserHandler(service user.UseCase, auth *Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bodyContent, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			log.Error().Err(err).Msg("could not read request body")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var req registerRequest
		if err = json.Unmarshal(bodyContent, &req); err != nil {
			log.Info().Err(err).Msg("invalid json")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}

		userID, err := service.CreateUser(r.Context(), req.Login, req.Password)
		if err != nil {
			if errors.Is(err, entity.ErrUserAlreadyExists) {
				http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
				return
			}
			if errors.Is(err, entity.ErrInvalidEntity) {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			log.Error().Err(err).Msg("error creating user")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		token, err := auth.GenerateToken(userID)
		if err != nil {
			log.Error().Err(err).Msg("error generating token")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		setJwtCookie(w, token)
		w.WriteHeader(http.StatusOK)
	}
}

func loginHandler(service user.UseCase, auth *Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bodyContent, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			log.Error().Err(err).Msg("could not read request body")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var req loginRequest
		if err = json.Unmarshal(bodyContent, &req); err != nil {
			log.Info().Err(err).Msg("invalid json")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}

		u, err := service.Login(r.Context(), req.Login, req.Password)
		if err != nil {
			if errors.Is(err, entity.ErrInvalidEntity) {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			if errors.Is(err, entity.ErrUserNotFound) || errors.Is(err, entity.ErrIncorrectPassword) {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			log.Error().Err(err).Msg("error authenticating user")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		token, err := auth.GenerateToken(u.ID)

		if err != nil {
			log.Error().Err(err).Msg("error generating token")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		setJwtCookie(w, token)
		w.WriteHeader(http.StatusOK)
	}
}

func MakeUserHandlers(r chi.Router, service user.UseCase, auth *Auth) {
	r.Post("/api/user/register", registerUserHandler(service, auth))
	r.Post("/api/user/login", loginHandler(service, auth))
}

func setJwtCookie(w http.ResponseWriter, token string) {
	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
}
