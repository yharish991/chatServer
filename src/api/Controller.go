package api

import "net/http"

// Controller interface for api
type Controller interface {
	Register()
	APIHandler(w http.ResponseWriter, r *http.Request)
	PostMessage(w http.ResponseWriter, r *http.Request)
	GetMessages(w http.ResponseWriter, r *http.Request)
}

