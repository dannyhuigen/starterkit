package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"log/slog"
	"os"
	"strconv"
	"sync"
)

var sqlxConnection *sqlx.DB
var pgxConnection *pgx.Conn
var mutex = sync.Mutex{}

func GetSqlXConnection() *sqlx.DB {
	if sqlxConnection == nil {
		initDbConnections()
	}
	return sqlxConnection
}

func GetPGXConnection() *pgx.Conn {
	if pgxConnection == nil {
		initDbConnections()
	}
	return pgxConnection
}

func initDbConnections() {
	mutex.Lock()
	defer mutex.Unlock()
	if sqlxConnection != nil {
		return
	}
	// Fetch environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPortStr := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSL_MODE")

	// Check which environment variables are missing and log them
	if dbHost == "" {
		log.Panic("DB_HOST environment variable is not set")
	}
	if dbPortStr == "" {
		log.Panic("DB_PORT environment variable is not set")
	}
	if dbUser == "" {
		log.Panic("DB_USER environment variable is not set")
	}
	if dbPassword == "" {
		log.Panic("DB_PASSWORD environment variable is not set")
	}
	if dbName == "" {
		log.Panic("DB_NAME environment variable is not set")
	}
	if dbSSLMode == "" {
		log.Panic("DB_SSL_MODE environment variable is not set")
	}

	// Cast port to integer
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Panicf("Invalid DB_PORT value: %s, it must be an integer", dbPortStr)
	}

	// Create connection string
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	// Using pgx as the driver for sqlx
	db, err := sqlx.Open("pgx", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	// Establish a connection using pgx directly (this is where pgx.Conn is returned)
	conn, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		log.Fatal(err)
	}

	sqlxConnection = db
	pgxConnection = conn
	HandleMigrations()
}

func HandleMigrations() {
	slog.Info("Handling migrations...")
	dbHost := os.Getenv("DB_HOST")
	dbPortStr := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSL_MODE")

	// Check which environment variables are missing and log them
	if dbHost == "" {
		log.Panic("DB_HOST environment variable is not set")
	}
	if dbPortStr == "" {
		log.Panic("DB_PORT environment variable is not set")
	}
	if dbUser == "" {
		log.Panic("DB_USER environment variable is not set")
	}
	if dbPassword == "" {
		log.Panic("DB_PASSWORD environment variable is not set")
	}
	if dbName == "" {
		log.Panic("DB_NAME environment variable is not set")
	}
	if dbSSLMode == "" {
		log.Panic("DB_SSL_MODE environment variable is not set")
	}

	m, err := migrate.New(
		"file://internal/database/migrations/",
		fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v", dbUser, dbPassword, dbHost, dbPortStr, dbName, dbSSLMode))
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Info("No migrations to apply")
			return
		}
		log.Fatal(err)
	}
}
