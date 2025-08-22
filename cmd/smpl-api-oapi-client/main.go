package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	api "github.com/vlladoff/micro-learn/internal/client"
)

func main() {
	apiClient := api.Client{
		Server: "http://localhost:8080",
		Client: &http.Client{},
	}

	resp, err := apiClient.GetPing(context.Background())
	if err != nil {
		log.Fatalf("Error calling GetPing: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf("Resonce code: %d\n", resp.StatusCode)

	var response map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Fatalf("Error decoding response: %v", err)
	}

	fmt.Printf("Response: %s\n", response["message"])
}
