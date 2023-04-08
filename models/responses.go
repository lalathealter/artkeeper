package models

type GetURLResponse struct {
	ID          string `field:"url_id"`
	Link        string `field:"url"`
	Description string `field:"url_description"`
	Username    string `field:"poster_id"`
}
