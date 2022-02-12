package user

import (
	"context"
	"github.com/thorgnir-go-study/yp-diploma/internal/entity"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) CreateUser(ctx context.Context, login string, password string) (entity.ID, error) {
	e, err := entity.NewUser(login, password)
	if err != nil {
		return entity.NilID, err
	}

	return s.repo.Create(ctx, *e)
}

func (s *Service) Login(ctx context.Context, login string, password string) (*entity.User, error) {
	user, err := s.repo.GetByLogin(ctx, login)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, entity.ErrUserNotFound
	}
	err = user.ValidatePassword(password)
	if err != nil {
		return nil, entity.ErrIncorrectPassword
	}

	return user, nil
}
