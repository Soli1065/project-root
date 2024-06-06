package main

import (
	"log"
	"net/http"
	"project-root/internal/api"
)

func main() {
	r := api.GatewayHandler()
	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", r)
}
