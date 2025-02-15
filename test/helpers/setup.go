package helpers

import (
	"os"
	"pack-management/internal/pkg/config"
	"pack-management/internal/pkg/database"

	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

func Setup() (*bun.DB, *fiber.App, func()) {
	os.Setenv("GO_ENV", "test")
	os.Setenv("BUNDEBUG", "2")

	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	app := fiber.New(fiber.Config{
		AppName: "test",
	})

	DBName := CreateDatabase(cfg)

	bunDB, err := database.NewDatabase(&database.Params{
		DBHost:     cfg.DBHost,
		DBPort:     cfg.DBPort,
		DBName:     DBName,
		DBUser:     cfg.DBUser,
		DBPassword: cfg.DBPassword,
	}).Connect()
	if err != nil {
		panic("connecting db: " + err.Error())
	}

	err = ExecMigrations(bunDB.DB, "migrations")
	if err != nil {
		panic("exec migrations: " + err.Error())
	}

	return bunDB, app, func() {
		app.Shutdown()
		ShutdownDB(bunDB.DB, DBName)
	}
}
