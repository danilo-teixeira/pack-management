package packevent

import (
	"context"
	"pack-management/internal/pkg/validator"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type (
	RepositoryParams struct {
		DB *bun.DB `validate:"required"`
	}

	mysqlRepository struct {
		db *bun.DB
	}
)

func NewMysqlRepository(params *RepositoryParams) Repository {
	params.validate()

	return &mysqlRepository{
		db: params.DB,
	}
}

func (p *RepositoryParams) validate() {
	err := validator.ValidateStruct(p)
	if err != nil {
		panic(err)
	}
}

func (r *mysqlRepository) Create(ctx context.Context, event *Entity) error {
	event.ID = r.newID()
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()

	_, err := r.db.NewInsert().Model(event.ToModel()).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *mysqlRepository) newID() string {
	return idPrefix + uuid.New().String()
}
