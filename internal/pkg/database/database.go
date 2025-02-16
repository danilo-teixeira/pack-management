package database

import (
	"database/sql"
	"fmt"
	"math"
	"runtime"

	"pack-management/internal/pkg/validator"

	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/extra/bundebug"
)

type (
	Service interface {
		Connect() (*bun.DB, error)
	}

	Params struct {
		DBHost         string `validate:"required"`
		DBPort         string `validate:"required"`
		DBName         string
		DBUser         string `validate:"required"`
		DBPassword     string `validate:"required"`
		ConnectionPool *ConnectionPool
	}

	ConnectionPool struct {
		MaxOpenConns    int
		IdleConnsFactor float64
	}

	service struct {
		dbHost          string
		dbPort          string
		dbName          string
		dbUser          string
		dbPassword      string
		MaxOpenConns    int
		IdleConnsFactor float64
	}
)

func WithPoolConfigLow() *ConnectionPool {
	return &ConnectionPool{
		MaxOpenConns:    2,
		IdleConnsFactor: 0.2,
	}
}

func WithPoolConfigHigh() *ConnectionPool {
	return &ConnectionPool{
		MaxOpenConns:    100,
		IdleConnsFactor: 0.2,
	}
}

func NewDatabase(params *Params) Service {
	validateParams(params)

	if params.ConnectionPool == nil {
		params.ConnectionPool = WithPoolConfigLow()
	}

	return &service{
		dbHost:          params.DBHost,
		dbPort:          params.DBPort,
		dbName:          params.DBName,
		dbUser:          params.DBUser,
		dbPassword:      params.DBPassword,
		IdleConnsFactor: params.ConnectionPool.IdleConnsFactor,
		MaxOpenConns:    params.ConnectionPool.MaxOpenConns,
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
		"%s:%s@tcp(%s:%s)/",
		s.dbUser,
		s.dbPassword,
		s.dbHost,
		s.dbPort,
	)

	if s.dbName != "" {
		host = fmt.Sprintf("%s%s?parseTime=true", host, s.dbName)
	}

	sqldb, err := sql.Open("mysql", host)
	if err != nil {
		return nil, err
	}

	s.setConnectionPool(sqldb)

	bunDB := bun.NewDB(sqldb, mysqldialect.New())
	bunDB.AddQueryHook(bundebug.NewQueryHook(bundebug.FromEnv("BUNDEBUG")))

	return bunDB, nil
}

func (s *service) setConnectionPool(db *sql.DB) {
	maxOpenConns := s.MaxOpenConns * runtime.GOMAXPROCS(0)
	maxIdleConns := int(math.RoundToEven(float64(maxOpenConns) * s.IdleConnsFactor))

	if maxIdleConns < 1 {
		maxIdleConns = 1
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
}
