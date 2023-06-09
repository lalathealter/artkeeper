package psql

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/lalathealter/artkeeper/config"
)

var currentDB *sql.DB

func init() {
	currentDB = connect()
	presetTables(currentDB)
}

func GetDB() *sql.DB {
	return currentDB
}

func connect() *sql.DB {
	db, err := sql.Open("postgres", config.Getnonempty("psqlconn"))
	if err != nil {
		log.Panicln(err)
	}

	err = db.Ping()
	if err != nil {
		log.Panicln(err)
	}

	fmt.Println("Database connected")
	return db
}

func presetTables(db *sql.DB) {
	initcommands := [...]string{
		`
			CREATE SCHEMA IF NOT EXISTS ak_data;
		`,
		`
			CREATE TABLE IF NOT EXISTS ak_data.users (
				user_name VARCHAR(36) NOT NULL,
				password_hash CHAR(120) NOT NULL,
				user_id SERIAL PRIMARY KEY NOT NULL,
				registration_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				UNIQUE(user_name)
			)
			;
		`,
		`
			CREATE TABLE IF NOT EXISTS ak_data.collections (
				collection_id SERIAL PRIMARY KEY NOT NULL,
				url_ids_collection INT [],
				collection_description TEXT,
				collection_tags VARCHAR(64) [],
				owner_id SERIAL NOT NULL,
				CONSTRAINT owner_id
					FOREIGN KEY(owner_id)
						REFERENCES ak_data.users(user_id)
						ON DELETE CASCADE
			)
			;
		`,
		`
			CREATE TABLE IF NOT EXISTS ak_data.urls (
				url TEXT NOT NULL,
				url_id BIGSERIAL PRIMARY KEY NOT NULL,
				url_description TEXT,
				poster_id SERIAL NOT NULL,
				CONSTRAINT poster_id
					FOREIGN KEY(poster_id)
						REFERENCES ak_data.users(user_id)
						ON DELETE CASCADE
			)
			;
		`,
	}

	for _, comm := range initcommands {
		_, err := db.Exec(comm)
		if err != nil {
			log.Panicln(err)
		}
	}
}
