package models

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func SetupSharding() {
	dsn := fmt.Sprintf(
		"host=%s user=%s dbname=%s port=%s %s",
		os.Getenv("DB_HOST"),
		"spqr-console",
		"spqr-console",
		os.Getenv("DB_PORT"),
		os.Getenv("EXTRA_PARAMS"),
	)
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %s", err))
	}
	defer func() {
		_ = conn.Close(context.Background())
	}()

	_, err = conn.Exec(context.Background(), "CREATE SHARDING RULE r1 COLUMNS id;")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "could not setup sharding: %s", err)
	}
	_, err = conn.Exec(context.Background(), "Add KEY RANGE krid1 FROM 1 TO 100 ROUTE TO sh1;")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "could not setup sharding: %s", err)
	}
	_, err = conn.Exec(context.Background(), "Add KEY RANGE krid2 FROM 100 TO 200 ROUTE TO sh2;")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "could not setup sharding: %s", err)
	}
}

func CreateTables() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s TimeZone=UTC %s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("EXTRA_PARAMS"),
	)

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %s", err))
	}
	defer func() {
		_ = conn.Close(context.Background())
	}()

	_, err = conn.Exec(context.Background(), "CREATE TABLE people (id bigserial, first_name text, last_name text, email text, PRIMARY KEY(id))")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "could not setup sharding: %s", err)
	}
}

func ConnectDatabase() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s TimeZone=UTC %s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("EXTRA_PARAMS"),
	)
	database, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database")
	}

	DB = database
}
