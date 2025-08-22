package main

import (
	"log"
	"net/http"

	"github.com/vlladoff/micro-learn/internal/app"
	api "github.com/vlladoff/micro-learn/internal/server"
)

func main() {
	server := app.NewSmplServer()
	handler := api.Handler(server)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
