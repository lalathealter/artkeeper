package models

type URLResponse struct {
	ID          string `field:"url_id"`
	Link        string `field:"url"`
	Description string `field:"url_description"`
	UserID      string `field:"poster_id"`
}

type CollectionResponse struct {
	ID string `field:"collection_id"`
	LinkIDs string `field:"url_ids_collection"`
	Tags string `field:"collection_tags"`
	Description string `field:"collection_description"`
	UserID string `field:"owner_id"`
}

type URLsFromCollectionResponse struct {
	LinkIDs string `field:"url_ids_collection"`
}
