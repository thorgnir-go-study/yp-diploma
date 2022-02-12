package user

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/thorgnir-go-study/yp-diploma/internal/entity"
	"github.com/thorgnir-go-study/yp-diploma/internal/usecase/user"
)

type PostgresUserRepository struct {
	dbpool *pgxpool.Pool
}

func NewPostgresUserRepository(dbpool *pgxpool.Pool) user.Repository {
	return &PostgresUserRepository{dbpool: dbpool}
}

func (p PostgresUserRepository) Get(ctx context.Context, id entity.ID) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresUserRepository) GetByLogin(ctx context.Context, login string) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresUserRepository) Create(ctx context.Context, e entity.User) (entity.ID, error) {
	//TODO implement me
	panic("implement me")
}
