package models

import (
	"database/sql"

	"github.com/lib/pq"
)


type GetURLRequest struct {
	ID     *ResourceID         `urlparam:"linkID"`
}
func (gr GetURLRequest) VerifyValues() error {
	return VerifyStruct(gr)
}

const selectOneURL = `
		SELECT * 
		FROM ak_data.urls 
		WHERE url_id=$1 
		;
	`
	
func (gr GetURLRequest) Call(db *sql.DB) (DBResult, error) {
	sqlstatement := selectOneURL
	sqlargs := []any{gr.ID}
	return db.Query(sqlstatement, sqlargs...)
}

type GetLatestURLsRequest struct {
	Offset *StringifiedInt `urlparam:"offset"`
	Limit  *StringifiedInt `urlparam:"limit"`
}

func (grLatest GetLatestURLsRequest) VerifyValues() error {
	return VerifyStruct(grLatest)
}
const selectAllUrls = ` 
		SELECT * 
		FROM ak_data.urls
		;
	`
const defaultPaginationLimit      = "10"
const selectLatestURLsWithPagination = `
		SELECT *
		FROM ak_data.urls
		ORDER BY url_id DESC
		LIMIT $1
		OFFSET $2
		;
	`
	
func (grLatest GetLatestURLsRequest) Call(db *sql.DB) (DBResult, error) {
	sqlstatement := selectLatestURLsWithPagination
	var sqlargs []any 
	if (*grLatest.Limit == "0") {
		sqlargs = []any{defaultPaginationLimit, grLatest.Offset}
	} else {
		sqlargs = []any{grLatest.Limit, grLatest.Offset}
	}
	return db.Query(sqlstatement, sqlargs...)
}


type DeleteURLRequest struct {
	// UserID *UserID `json:"userID"`
	LinkID *ResourceID `json:"linkID"`
}

func (dr DeleteURLRequest) VerifyValues() error {
	return VerifyStruct(dr)
}

const deleteOneURL = `
		DELETE FROM ak_data.urls
		WHERE url_id=$1
		;
	`

func (dr DeleteURLRequest) Call(db *sql.DB) (DBResult, error) {
	sqlstatement := deleteOneURL
	sqlargs := []any{ dr.LinkID }
	return db.Exec(sqlstatement, sqlargs...)
}

type PostURLRequest struct {
	Link        *InputLink   `json:"link"`
	Description *Description `json:"description"`
	UserID      *UserID      `json:"userID"`
}

func (pr PostURLRequest) VerifyValues() error {
	return VerifyStruct(pr)
}

const insertOneURL = `
	INSERT INTO ak_data.urls(url, url_description, poster_id) 
	VALUES($1, $2, $3)
	;
`
	
func (pr PostURLRequest) Call(db *sql.DB) (DBResult, error) {
	sqlstatement := insertOneURL
	sqlargs := []any{ pr.Link, pr.Description, pr.UserID }
	return db.Exec(sqlstatement, sqlargs...)
}

type PostCollectionRequest struct {
	LinkIDs     []*ResourceID    `json:"linkIDs"`
	Description *Description `json:"description"`
	UserID      *UserID      `json:"userID"`
}

func (pcr PostCollectionRequest) VerifyValues() error {
	return VerifyStruct(pcr)
}

const insertOneCollection = `
		INSERT INTO ak_data.collections(url_ids_collection, owner_id, collection_description)
		VALUES($1, $2, $3)
		;
	`
func (pcr PostCollectionRequest) Call(db *sql.DB) (DBResult, error) {
	sqlstatement := insertOneCollection
	sqlargs := []any{
		pq.Array(pcr.LinkIDs),
		pcr.UserID,
		pcr.Description,
	}
	return db.Query(sqlstatement, sqlargs...)
}

type PutInCollectionRequest struct {
	LinkID *ResourceID `json:"linkID"`
	CollID *ResourceID `json:"collID"`
}

func (putcr PutInCollectionRequest) VerifyValues() error {
	return VerifyStruct(putcr)
}

const updateLinksInCollection = `
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
	
func (putcr PutInCollectionRequest) Call(db *sql.DB) (DBResult, error) {
	sqlstatement := updateLinksInCollection
	sqlargs := []any{
		putcr.LinkID,
		putcr.CollID,
	}
	return db.Query(sqlstatement, sqlargs...)
}

type GetCollectionRequest struct {
	ID *ResourceID `urlparam:"collID"`
}

func (gcr GetCollectionRequest) VerifyValues() error {
	return VerifyStruct(gcr)
}

const selectOneCollection = `
		SELECT *
		FROM ak_data.collections
		WHERE collection_id=$1
		;
	`
func (gcr GetCollectionRequest) Call(db *sql.DB) (DBResult, error) {

	sqlstatement := selectOneCollection
	sqlargs := []any{ gcr.ID }

	return db.Query(sqlstatement, sqlargs...)
}

type DeleteCollectionRequest struct {
	CollID *ResourceID `json:"collID"`
}

func (dcr DeleteCollectionRequest) VerifyValues() error {
	return VerifyStruct(dcr)
}

const deleteOneCollection = `
	DELETE FROM ak_data.collections 
	WHERE collection_id=$1
	;
`

func (dcr DeleteCollectionRequest) Call(db *sql.DB) (DBResult, error) {
	sqlstatement := deleteOneCollection 
	sqlargs := []any{ dcr.CollID }
	return db.Exec(sqlstatement, sqlargs...)
}


