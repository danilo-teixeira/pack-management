package holiday

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

func (r *mysqlRepository) Create(ctx context.Context, holiday *Entity) error {
	holiday.ID = r.newID()
	holiday.CreatedAt = time.Now()
	holiday.UpdatedAt = time.Now()

	_, err := r.db.NewInsert().Model(holiday.ToModel()).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *mysqlRepository) BulkCreate(ctx context.Context, holidays []*Entity) error {
	holidayModels := make([]*Model, 0, len(holidays))

	for _, holiday := range holidays {
		holiday.ID = r.newID()
		holiday.CreatedAt = time.Now()
		holiday.UpdatedAt = time.Now()

		holidayModels = append(holidayModels, holiday.ToModel())
	}

	_, err := r.db.NewInsert().Model(&holidayModels).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *mysqlRepository) ListByYear(ctx context.Context, year string) ([]*Entity, error) {
	holidays := make([]*Model, 0)
	dateGte := year + "-01-01"
	dateLte := year + "-12-31"

	err := r.db.NewSelect().
		Model(&holidays).
		Where("date BETWEEN ? and ?", dateGte, dateLte).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	entities := []*Entity{}
	for _, holiday := range holidays {
		entities = append(entities, holiday.ToEntity())
	}

	return entities, nil
}

func (r *mysqlRepository) newID() string {
	return idPrefix + uuid.New().String()
}
