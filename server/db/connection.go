package db

import (
	"context"
	"fmt"
	"log"

	"pfserver/config"

	"github.com/jackc/pgx/v4/pgxpool"
)

var conn *pgxpool.Pool

func init() {

	// check if the connection is already there
	if conn != nil {
		return
	}

	var dbErr error
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	DB_URL := config.GetEnv("POSTGRES_URL")
	// if we are in testing mode
	if config.IsTestMode() {
		user := config.GetEnv("POSTGRES_USER")
		pass := config.GetEnv("POSTGRES_PASSWORD")
		db := config.GetEnv("POSTGRES_DB")
		DB_URL = fmt.Sprintf(
			"postgres://%s:%s@localhost:5432/%s",
			user, pass, db,
		)
	}
	conn, dbErr = pgxpool.Connect(context.Background(), DB_URL)

	if dbErr != nil {
		log.Fatal("Unable to connect to database, ERROR: ", dbErr)
	}

	//? Migrate the database
	if !config.IsTestMode() {
		Migrate(DB_URL)
	}
}

func Conn() *pgxpool.Pool {
	return conn
}
