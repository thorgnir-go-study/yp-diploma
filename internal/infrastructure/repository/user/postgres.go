package user

import (
	"context"
	"errors"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/thorgnir-go-study/yp-diploma/internal/entity"
	"github.com/thorgnir-go-study/yp-diploma/internal/usecase/user"
)

type PostgresUserRepository struct {
	dbpool *pgxpool.Pool
}

type userDbEntity struct {
	ID       uuid.UUID `db:"id"`
	Login    string    `db:"login"`
	Password string    `db:"password"`
}

func NewPostgresUserRepository(dbpool *pgxpool.Pool) user.Repository {
	return &PostgresUserRepository{dbpool: dbpool}
}

func (p PostgresUserRepository) GetByLogin(ctx context.Context, login string) (*entity.User, error) {
	var dbEntity userDbEntity
	if err := pgxscan.Get(ctx, p.dbpool, &dbEntity, `SELECT id, login, password
FROM gophermart.users
WHERE login=$1`, login); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrUserNotFound
		}
		return nil, err
	}

	return &entity.User{
		ID:       dbEntity.ID,
		Login:    dbEntity.Login,
		Password: dbEntity.Password,
	}, nil
}

func (p PostgresUserRepository) Create(ctx context.Context, e entity.User) (entity.ID, error) {
	var userID uuid.UUID
	err := p.dbpool.QueryRow(ctx, `INSERT INTO gophermart.users (login, password) 
VALUES ($1, $2) 
RETURNING id`, e.Login, e.Password).Scan(&userID)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "login_unique" {
				return entity.NilID, entity.ErrUserAlreadyExists
			}
		}
		return entity.NilID, err
	}
	return userID, nil
}
