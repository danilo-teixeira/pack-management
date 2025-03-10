package person

import (
	"context"
	"database/sql"
	"errors"
	"pack-management/internal/pkg/validator"
	"time"

	"pack-management/internal/pkg/uuid"

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

func (r *mysqlRepository) Create(ctx context.Context, person *Entity) error {
	person.ID = r.newID()
	person.CreatedAt = time.Now()
	person.UpdatedAt = time.Now()

	_, err := r.db.NewInsert().Model(person.ToModel()).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *mysqlRepository) GetByName(ctx context.Context, name string) (*Entity, error) {
	person := Model{}

	err := r.db.NewSelect().Model(&person).Where("name = ?", name).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return person.ToEntity(), nil
}

func (r *mysqlRepository) newID() string {
	return idPrefix + uuid.New().String()
}
