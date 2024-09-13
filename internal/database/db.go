package database

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/config"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/db"
)

// Reusable SQLC 
var DB *db.Queries

func InitDB(dbConfig config.DatabaseConfig) {
	var err error
	uri := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.DBHost,
		dbConfig.DBPort,
		dbConfig.DBUserName,
		dbConfig.DBUserPassword,
		dbConfig.DBName,
	)
	conn, err := pgx.Connect(context.Background(), uri)
	if err != nil {
		slog.Error("Failed to connect to database with error: " + err.Error())
		slog.Info("The connection string used was: " + uri)
		os.Exit(1)
	}
	DB = db.New(conn)
}
