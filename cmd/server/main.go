package main

import (
	"net/http"

	handler "github.com/lebendig13/metrics/internal/handler"
	models "github.com/lebendig13/metrics/internal/model"
)

func main() {
	memStorage := models.NewMemStorage()
	server := handler.NewServer(memStorage)

	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, server.UpdateMetricHandler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
