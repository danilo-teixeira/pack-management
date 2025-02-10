package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"

	"pack-management/internal/pkg/validator"
)

type (
	Service interface {
		Connect() (*bun.DB, error)
	}

	Params struct {
		DBHost     string `validate:"required"`
		DBPort     string `validate:"required"`
		DBName     string `validate:"required"`
		DBUser     string `validate:"required"`
		DBPassword string `validate:"required"`
	}

	service struct {
		dbHost     string
		dbPort     string
		dbName     string
		dbUser     string
		dbPassword string
	}
)

// TODO: impement connection pool
func NewDatabase(params *Params) Service {
	validateParams(params)

	return &service{
		dbHost:     params.DBHost,
		dbPort:     params.DBPort,
		dbName:     params.DBName,
		dbUser:     params.DBUser,
		dbPassword: params.DBPassword,
	}
}

func validateParams(params *Params) {
	err := validator.ValidateStruct(params)
	if err != nil {
		panic(err)
	}
}

func (s *service) Connect() (*bun.DB, error) {
	host := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		s.dbUser,
		s.dbPassword,
		s.dbHost,
		s.dbPort,
		s.dbName,
	)

	sqldb, err := sql.Open("mysql", host)
	if err != nil {
		return nil, err
	}

	return bun.NewDB(sqldb, mysqldialect.New()), nil
}
