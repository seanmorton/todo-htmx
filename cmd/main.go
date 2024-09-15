package main

import (
	"database/sql"
	"embed"
	"log/slog"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/seanmorton/todo-htmx/internal/data"
	"github.com/seanmorton/todo-htmx/internal/handlers"
)

//go:embed public
var publicDir embed.FS

func main() {
	tz, err := time.LoadLocation("America/Chicago")
	if err != nil {
		slog.Error("failed loading timezone location", "err", err)
		os.Exit(1)
	}

	dbFile := os.Getenv("DB_FILE")
	if dbFile == "" {
		slog.Error("DB_FILE env not set")
		os.Exit(1)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	dbConn, err := sql.Open("sqlite3", dbFile)
	defer dbConn.Close()
	if err != nil {
		slog.Error("failed opening db", "err", err)
		os.Exit(1)
	}
	dbConn.Exec("PRAGMA foreign_keys = ON;")

	db := data.NewDB(dbConn)
	server := handlers.NewServer(db, tz)

	slog.Info("starting server...")
	err = server.Start(port, publicDir)
	if err != nil {
		slog.Error("failed to start server", "err", err)
		os.Exit(1)
	}
}
