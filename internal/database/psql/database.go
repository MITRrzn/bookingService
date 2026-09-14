package psql

import (
	"database/sql"
	"log"
	"os"
)

func PsqlConnect() (*sql.DB, error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	if pingErr := db.Ping(); pingErr != nil {
		closeErr := db.Close()
		if closeErr != nil {
			log.Printf("failed to close postgres connection: %v\n", closeErr)
		}
		return nil, pingErr
	}

	return db, nil
}
