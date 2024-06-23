package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/seanmorton/todo-htmx/internal/app"
)

func main() {
	tz, _ := time.LoadLocation("America/Chicago") // TODO configure

	dbFile := os.Getenv("DB_FILE")
	if dbFile == "" {
		slog.Error("DB_FILE env not set")
		os.Exit(1)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	db, err := sql.Open("sqlite3", dbFile)
	defer db.Close()
	if err != nil {
		slog.Error("failed opening db", "err", err)
	}
	db.Exec("PRAGMA foreign_keys = ON;")

	tasksDB := app.NewTasksDB(db)
	server := app.NewServer(tasksDB, tz)

	slog.Info("starting server")
	server.RegisterRoutes()
	err = http.ListenAndServe(port, nil)
	if err != nil {
		slog.Error("failed to start server", "err", err)
	}
}
