package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"

	"ydx-goadv-gophkeeper/internal/logger"
	"ydx-goadv-gophkeeper/internal/model"
	"ydx-goadv-gophkeeper/internal/model/consts"
	"ydx-goadv-gophkeeper/internal/model/errs"
)

type UserRepository interface {
	Register(context.Context, model.User) (int32, error)
	Login(context.Context, model.User) (int32, error)
}

type userRepository struct {
	log *zap.SugaredLogger
	db  DBProvider
}

func NewUserRepository(db DBProvider) UserRepository {
	return &userRepository{log: logger.NewLogger("auth-repo"), db: db}
}

func (repo *userRepository) Register(ctx context.Context, user model.User) (int32, error) {
	conn, err := repo.db.GetConnection(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Release()

	queryRow := conn.QueryRow(ctx, "insert into users (username, password) values ($1, $2) returning id", user.Username, user.Password)
	var userId int32
	err = queryRow.Scan(&userId)
	if pgError, ok := err.(*pgconn.PgError); ok && pgError.Code == consts.UniqueViolation {
		return 0, errs.ErrUserAlreadyExist
	}
	if err != nil {
		repo.log.Errorf("failed to save user '%s': %v", user.Username, err)
		return 0, fmt.Errorf("failed to save user '%s': %v", user.Username, err)
	}
	return userId, nil
}

func (repo *userRepository) Login(ctx context.Context, user model.User) (int32, error) {
	conn, err := repo.db.GetConnection(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Release()

	queryRow := conn.QueryRow(ctx, "select id from users where username = $1 and password = $2", user.Username, user.Password)
	var userId int32
	err = queryRow.Scan(&userId)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, errs.ErrUserNotFound
	}
	if err != nil {
		repo.log.Errorf("failed to get user '%s': %v", user.Username, err)
		return 0, fmt.Errorf("failed to get user '%s': %v", user.Username, err)
	}

	return userId, nil
}
