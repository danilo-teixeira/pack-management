package helpers

import (
	"database/sql"
	"fmt"
	"pack-management/internal/pkg/config"
	"pack-management/internal/pkg/database"
	"pack-management/internal/pkg/helpers"
	"pack-management/internal/pkg/uuid"
	"strings"

	migrate "github.com/rubenv/sql-migrate"
)

func CreateDatabase(cfg *config.Config) string {

	DBName := "it_" + uuid.New().String()
	DBName = strings.ReplaceAll(DBName, "-", "_")

	bunDB, err := database.NewDatabase(&database.Params{
		DBHost:     cfg.DBHost,
		DBPort:     cfg.DBPort,
		DBUser:     cfg.DBUser,
		DBPassword: cfg.DBPassword,
	}).Connect()
	if err != nil {
		panic(err)
	}
	defer bunDB.Close()

	_, err = bunDB.DB.Exec(`CREATE DATABASE ` + DBName + `;`)
	if err != nil {
		panic(fmt.Sprintf("error creating db: %s", err.Error()))
	}

	return DBName
}

func ShutdownDB(db *sql.DB, dbName string) {
	_, err := db.Exec(`DROP DATABASE ` + dbName + `;`)
	if err != nil {
		panic(fmt.Sprintf("error dropping db: %s", err.Error()))
	}

	db.Close()
}

func ExecMigrations(database *sql.DB, table string) error {
	rootDirectory, err := helpers.GetRootDirectory()
	if err != nil {
		return err
	}

	seeders := &migrate.FileMigrationSource{
		Dir: rootDirectory + "scripts/db/migrations",
	}

	migrate.SetTable(table)
	_, err = migrate.Exec(database, "mysql", seeders, migrate.Up)

	if err != nil {
		return err
	}

	return nil
}
