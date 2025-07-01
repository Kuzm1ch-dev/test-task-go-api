package http

import (
	"github.com/gorilla/mux"
)

func SetupRoutes(handler *Handler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/tasks", handler.CreateTask).Methods("POST")
	r.HandleFunc("/tasks/{id}", handler.GetTask).Methods("GET")
	r.HandleFunc("/tasks/{id}", handler.DeleteTask).Methods("DELETE")

	return r
}
