package database

import "github.com/jmoiron/sqlx"

func ViewTables(db *sqlx.DB) ([]string, error) {
	var tables []string
	err := db.Select(&tables, "SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		return tables, err
	}
	return tables, nil
}
