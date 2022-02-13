package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/thorgnir-go-study/yp-diploma/internal/api/middleware/jwt"
	"github.com/thorgnir-go-study/yp-diploma/internal/entity"
	"github.com/thorgnir-go-study/yp-diploma/internal/pkg/auth"
)

type Auth struct {
	jwtWrapper *auth.JwtWrapper
}

func NewAuth(jwtWrapper *auth.JwtWrapper) *Auth {
	return &Auth{jwtWrapper: jwtWrapper}
}

func (s *Auth) GenerateToken(userID entity.ID) (string, error) {
	return s.jwtWrapper.GenerateToken(userID)
}

func (s *Auth) RegisterAuthMiddleware(r chi.Router) {
	r.Use(jwt.AuthMiddleware(s.jwtWrapper))
}
