package models

import (
	"database/sql"

	"github.com/lib/pq"
)


type GetURLRequest struct {
	ID     *ResourceID         `urlparam:"0"`
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
	Offset *StringifiedInt `urlquery:"offset"`
	Limit  *StringifiedInt `urlquery:"limit"`
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
	LinkID *ResourceID `urlparam:"0"`
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
		VALUES((
			SELECT ARRAY (
				SELECT url_id 
				FROM ak_data.urls
				WHERE url_id IN (
					SELECT unnest($1::INT[])
				)
			)
		), $2, $3)
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
	CollID *ResourceID `urlparam:"1"`
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
		RETURNING collection_id
		;
	`

const checkIfLinkExists = `
	SELECT EXISTS(
		SELECT 1 
		FROM ak_data.urls
		WHERE url_id=$1
	)
	;
`

func (putcr PutInCollectionRequest) Call(db *sql.DB) (DBResult, error) {
	doesLinkExist := false
	err := db.QueryRow(checkIfLinkExists, putcr.LinkID).Scan(&doesLinkExist)
	if !doesLinkExist {
		return nil, err
	}

	sqlstatement := updateLinksInCollection
	sqlargs := []any{
		putcr.LinkID,
		putcr.CollID,
	}
	return db.Exec(sqlstatement, sqlargs...)
}

type GetCollectionRequest struct {
	ID *ResourceID `urlparam:"0"`
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
	CollID *ResourceID `urlparam:"0"`
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

type DeleteURLFromCollectionRequest struct {
	CollID *ResourceID `urlparam:"2"`
	LinkID *ResourceID `urlparam:"0"`
}

func (dlcr DeleteURLFromCollectionRequest) VerifyValues() error {
	return VerifyStruct(dlcr)
}

const deleteURLFromCollection = `
		UPDATE ak_data.collections
		SET url_ids_collection = (
			SELECT ARRAY (
				SELECT unnest(url_ids_collection) 
				EXCEPT SELECT $1
			)
		)
		WHERE collection_id=$2
		;
	`
	

func (dlcr DeleteURLFromCollectionRequest) Call (db *sql.DB) (DBResult, error) {
	sqlstatement := deleteURLFromCollection
	sqlargs := []any{  dlcr.LinkID, dlcr.CollID }
	return db.Exec(sqlstatement, sqlargs...)
}

type GetURLsFromCollectionRequest struct {
	ID *ResourceID `urlparam:"1"`
}

const selectURLsFromCollection = `
	SELECT *
	FROM ak_data.urls 
	WHERE url_id IN (
		SELECT unnest(url_ids_collection)
		FROM ak_data.collections 
		WHERE collection_id=$1
	)
	;
`

func (glcr GetURLsFromCollectionRequest) VerifyValues() (error) {
	return VerifyStruct(glcr) 
}

func (glcr GetURLsFromCollectionRequest) Call(db *sql.DB) (DBResult, error) {
	sqlstatement := selectURLsFromCollection 
	sqlargs := []any{ glcr.ID }
	return db.Query(sqlstatement, sqlargs...)
}

type AttachTagToCollectionRequest struct {
	TagName *Tag `urlparam:"0"`  
	CollID *ResourceID `urlparam:"2"`
}

func (attag AttachTagToCollectionRequest) VerifyValues() (error) {
	return VerifyStruct(attag)
}

const updateTagsInCollection = `
		UPDATE ak_data.collections
		SET collection_tags = (
			SELECT ARRAY (
				SELECT DISTINCT * 
				FROM unnest(
					array_append(collection_tags, $1)
				)
			)
		)
		WHERE collection_id=$2
		RETURNING collection_id
		;
`

func (attag AttachTagToCollectionRequest) Call(db *sql.DB) (DBResult, error) {
	sqlstatement := updateTagsInCollection
	sqlargs := []any{ attag.TagName, attag.CollID }
	return db.Exec(sqlstatement, sqlargs...)
}

type DetachTagFromCollectionRequest struct {
	TagName *Tag `urlparam:"0"`
	CollID *ResourceID `urlparam:"2"`
}

func (detag DetachTagFromCollectionRequest) VerifyValues() error {
	return VerifyStruct(detag)
}

const deleteTagFromCollection = `
		UPDATE ak_data.collections
		SET collection_tags = (
			SELECT ARRAY (
				SELECT unnest(collection_tags) 
				EXCEPT SELECT $1
			)
		)
		WHERE collection_id=$2
		RETURNING collection_id
		;
	`


func (detag DetachTagFromCollectionRequest) Call(db *sql.DB) (DBResult, error) {
	sqlstatement := deleteTagFromCollection
	sqlargs := []any{ detag.TagName, detag.CollID }
	return db.Exec(sqlstatement, sqlargs...)
}
