package migrator

import (
	"log/slog"
	"os"

	"github.com/Grubiha/auth_session/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	slog.Debug("config loaded", "config", cfg)

	databaseUrl := config.PgUrl(cfg)
	slog.Debug("connecting to database", "url", databaseUrl)

	m, err := migrate.New("file://migrations", databaseUrl)
	if err != nil {
		slog.Error("Failed to create migrate instance", "error", err.Error())
		os.Exit(1)
	}

	err = m.Up()
	if err != nil {
		slog.Error("Failed to migrate database", "error", err.Error())
	} else {
		slog.Info("migrated database")
	}
}
