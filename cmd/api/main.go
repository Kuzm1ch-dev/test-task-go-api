package main

import (
	"api/internal/services/task"
	"api/internal/storage/memory"
	h "api/internal/transport/http"
	"log"
	"net/http"
)

func main() {
	storage := memory.New()
	taskService := task.New(storage)
	handler := h.NewHandler(taskService)

	router := h.SetupRoutes(handler)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
