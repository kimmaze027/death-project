package main

import (
	"log"
	"net/http"

	"github.com/kimmaze027/death-project/apps/backend/internal/api"
	"github.com/kimmaze027/death-project/apps/backend/internal/store"
)

func main() {
	memoryStore := store.NewMemoryStore()
	server := api.NewServer(memoryStore)

	addr := ":8080"
	log.Printf("backend server listening on %s", addr)
	if err := http.ListenAndServe(addr, server.Handler()); err != nil {
		log.Fatal(err)
	}
}
