package migrations

import (
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func getPathMig() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Ошибка в директории", err)
	}
	return "file://" + filepath.Join(wd, "migrations")
}

func RunMigration(dsn string) {
	mig, err := migrate.New(
		getPathMig(),
		dsn,
	)
	if err != nil {
		log.Fatal("ошибка миграций", err)
	}

	if err := mig.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Ошибка при миграции", err)
	}

	log.Println("Миграция успешна")

}
