package pack

import (
	"context"
	"database/sql"
	"errors"
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

func (r *mysqlRepository) Create(ctx context.Context, pack *Entity) error {
	pack.ID = r.newID()
	pack.CreatedAt = time.Now()
	pack.UpdatedAt = time.Now()

	_, err := r.db.NewInsert().Model(pack.ToModel()).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *mysqlRepository) UpdateByID(ctx context.Context, ID string, pack *Entity) error {
	pack.UpdatedAt = time.Now()

	_, err := r.db.NewUpdate().Model(pack.ToModel()).Where("id = ?", ID).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *mysqlRepository) GetByID(ctx context.Context, ID string, withEvents bool) (*Entity, error) {
	pack := Model{}

	query := r.db.NewSelect().
		Model(&pack).
		Where("pack.id = ?", ID).
		Limit(1).
		Relation("Receiver").
		Relation("Sender")

	if withEvents {
		query.Relation("Events")
	}

	if err := query.Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return pack.ToEntity(), nil
}

func (r *mysqlRepository) newID() string {
	return idPrefix + uuid.New().String()
}
