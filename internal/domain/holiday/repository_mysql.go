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

func (r *mysqlRepository) GetByDate(ctx context.Context, date string) (*Entity, error) {
	holiday := Model{}

	err := r.db.NewSelect().Model(&holiday).Where("date = ?", date).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return holiday.ToEntity(), nil
}

func (r *mysqlRepository) newID() string {
	return idPrefix + uuid.New().String()
}
