package user

import (
	"context"
	"github.com/thorgnir-go-study/yp-diploma/internal/entity"
)

type Reader interface {
	GetByLogin(ctx context.Context, login string) (*entity.User, error)
}

type Writer interface {
	Create(ctx context.Context, e entity.User) (entity.ID, error)
}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	CreateUser(ctx context.Context, login string, password string) (entity.ID, error)
	Login(ctx context.Context, login string, password string) (*entity.User, error)
}
