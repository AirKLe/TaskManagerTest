package main

import (
	"TaskManager/internal/api"
	"TaskManager/internal/service"
	"TaskManager/internal/storage"
	"log"
	"net/http"
)

func main() {
	storе := storage.NewInMemoryTaskStorage()
	svc := service.NewTaskService(storе)
	handler := api.NewTaskHandler(svc)

	http.Handle("/tasks/", handler)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
