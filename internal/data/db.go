package data

import (
	"database/sql"
)

type DB struct {
	dbConn *sql.DB
}

func NewDB(dbConn *sql.DB) DB {
	return DB{dbConn}
}
