package repository

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type PostgreRepository struct {
}

func (p *PostgreRepository) CreateClientDatabase() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatalln("Error creating database connection: ", err)
	}
	return db
}
