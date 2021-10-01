package postgres

import "database/sql"

type Store struct {
	DB *sql.DB
}
