package api

import (
	"database/sql"
)

type ApiServer struct {
	addr string
	db   *sql.DB
}

func NewServier(addr string, db *sql.DB) *ApiServer {
	return &ApiServer{
		addr: addr,
		db:   db,
	}
}
