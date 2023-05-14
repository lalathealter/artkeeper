package controllers

import "net/http"

func HelpURLHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("#documentation for interacting with urls"))
}

func HelpCollectionHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("#documentation for interacting with collections"))
}
