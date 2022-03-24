package db

import (
	"fmt"
	"log"
	"path/filepath"

	"pfserver/utils"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/golang-migrate/migrate/v4"
)

func Migrate(DB_URL string) {
	// get migrations source, it must start with "file://"
	migrationSource := filepath.Join("file://", utils.RootPath(), "db", "migrations")
	fmt.Println("Migration Source: ", migrationSource)

	migration, migErr := migrate.New(migrationSource, DB_URL)

	if migErr != nil {
		log.Println("migrate.New: ", migErr.Error())
	}
	migV, migVDitry, migVErr := migration.Version()
	if migVDitry && migVErr == nil {
		migration.Force(int(migV) - 1)
	}

	migUpErr := migration.Up()

	if migUpErr != nil {
		log.Println("migrate.Up: ", migUpErr)
	}
}
