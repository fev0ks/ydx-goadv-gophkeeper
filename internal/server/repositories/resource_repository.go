package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"

	"ydx-goadv-gophkeeper/internal/server/model"
	"ydx-goadv-gophkeeper/internal/server/model/errs"
	"ydx-goadv-gophkeeper/pkg/logger"
	"ydx-goadv-gophkeeper/pkg/model/enum"
)

type ResourceRepository interface {
	Save(ctx context.Context, resource *model.Resource) error
	Get(ctx context.Context, resId int32, userId int32) (*model.Resource, error)
	GetResDescriptionsByType(ctx context.Context, userId int32, resType enum.ResourceType) ([]*model.ResourceDescription, error)
	Delete(ctx context.Context, resId int32, userId int32) error
}

type resourceRepository struct {
	log *zap.SugaredLogger
	db  DBProvider
}

func NewResourceRepository(db DBProvider) ResourceRepository {
	return &resourceRepository{log: logger.NewLogger("res-repo"), db: db}
}

func (s *resourceRepository) Save(ctx context.Context, resource *model.Resource) error {
	conn, err := s.db.GetConnection(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	var resId int32
	row := conn.QueryRow(
		ctx,
		"insert into resources(user_id, type, data, meta) values ($1, $2, $3, $4) RETURNING id",
		resource.UserId,
		resource.Type,
		resource.Data,
		resource.Meta,
	)
	err = row.Scan(&resId)
	if err != nil {
		return err
	}
	resource.Id = resId

	return err
}

func (r *resourceRepository) Get(ctx context.Context, resId int32, userId int32) (*model.Resource, error) {
	var result model.Resource
	conn, err := r.db.GetConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	var row pgx.Row
	row = conn.QueryRow(ctx, "select id, user_id, type, meta, data from resources where id = $1 and user_id = $2", resId, userId)
	err = row.Scan(&result.Id, &result.UserId, &result.Type, &result.Meta, &result.Data)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errs.ErrResNotFound
	}
	if err != nil {
		r.log.Errorf("failed to parse get resourse '%d' result row: %v", resId, err)
		return nil, err
	}
	return &result, err
}

func (s *resourceRepository) GetResDescriptionsByType(ctx context.Context, userId int32, resType enum.ResourceType) ([]*model.ResourceDescription, error) {
	var results []*model.ResourceDescription
	conn, err := s.db.GetConnection(ctx)
	if err != nil {
		return results, err
	}
	defer conn.Release()

	var rows pgx.Rows
	if resType == enum.Nan {
		rows, err = conn.Query(
			ctx,
			"select id, meta, type from resources where user_id = $1",
			userId,
		)
	} else {
		rows, err = conn.Query(
			ctx,
			"select id, meta, type from resources where user_id = $1 and type = $2",
			userId,
			resType,
		)
	}
	defer rows.Close()
	for rows.Next() {
		resDescr := &model.ResourceDescription{}
		err := rows.Scan(&resDescr.Id, &resDescr.Meta, &resDescr.Type)
		if err != nil {
			s.log.Errorf("failed to read '%d' resources of userId '%d': %v", resType, userId, err)
			return nil, fmt.Errorf("failed to read '%d' resources of userId '%d': %v", resType, userId, err)
		}
		results = append(results, resDescr)
	}
	return results, err
}

func (s *resourceRepository) Delete(ctx context.Context, resId int32, userId int32) error {
	conn, err := s.db.GetConnection(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, "delete from resources where id = $1 and user_id = $2", resId, userId)
	return err
}
