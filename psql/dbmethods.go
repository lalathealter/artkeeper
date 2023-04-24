package psql

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/lalathealter/artkeeper/config"
)

var (
	InsertOneURL string = `
		INSERT INTO ak_data.urls(url, url_description, poster_id) 
		VALUES($1, $2, $3)
		;
	`
	SelectOneURL string = `
		SELECT * 
		FROM ak_data.urls 
		WHERE url_id=$1 
		;
	`
	SelectAllURLs string = `
		SELECT * 
		FROM ak_data.urls
		;
	`
	DefaultPaginationLimit      string = "10"
	SelectLatestURLsWithPagination string = `
		SELECT *
		FROM ak_data.urls
		ORDER BY url_id DESC
		LIMIT $1
		OFFSET $2
		;
	`
	DeleteOneURL = `
		DELETE FROM ak_data.urls
		WHERE url_id=$1
		;
	`

	// ====================

	InsertOneCollection string = `
		INSERT INTO ak_data.collections(url_ids_collection, owner_id, collection_description)
		VALUES($1, $2, $3)
		;
	`
	UpdateLinksInCollection string = `
		UPDATE ak_data.collections
		SET url_ids_collection = (
			SELECT ARRAY (
				SELECT DISTINCT * 
				FROM unnest(
					array_append(url_ids_collection, $1)
				)
			)
		)
		WHERE collection_id=$2
		;
	`

	SelectOneCollection string = `
		SELECT *
		FROM ak_data.collections
		WHERE collection_id=$1
		;
	`
)

var DB *sql.DB

func init() {
	DB = connect()
	presetTables(DB)
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
				user_id SERIAL PRIMARY KEY NOT NULL,
				registration_date TIMESTAMPTZ NOT NULL DEFAULT NOW() 
			)
			;
		`,
		`
			CREATE TABLE IF NOT EXISTS ak_data.collections (
				collection_id SERIAL PRIMARY KEY NOT NULL,
				url_ids_collection INT [],
				collection_description TEXT,
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
