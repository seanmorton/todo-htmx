package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/seanmorton/todo-htmx/internal/app"
	"github.com/seanmorton/todo-htmx/internal/data"
)

func main() {
	tz, _ := time.LoadLocation("America/Chicago") // TODO configure or use client tz

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
	}
	dbConn.Exec("PRAGMA foreign_keys = ON;")

	db := data.NewDB(dbConn)
	server := app.NewServer(db, tz)

	slog.Info("starting server")
	server.RegisterRoutes()
	err = http.ListenAndServe(port, nil)
	if err != nil {
		slog.Error("failed to start server", "err", err)
	}
}
