package main

import (
	"log"
	"net/http"
)


func NewHandler() http.Handler {
	mux := http.NewServeMux()
	store := NewMemoryStore()

	// Routing B1: Sumber Daya
	mux.HandleFunc("POST /sumber-daya", store.handlePostSumberDaya)
	mux.HandleFunc("GET /sumber-daya", store.handleGetSumberDaya)
	mux.HandleFunc("GET /sumber-daya/{id}", store.handleGetSumberDayaByID)

	// Routing B2 & B3, B6: Peminjaman
	mux.HandleFunc("POST /peminjaman", store.handlePostPeminjaman)
	mux.HandleFunc("GET /peminjaman", store.handleGetPeminjaman)
	mux.HandleFunc("GET /peminjaman/{id}", store.handleGetPeminjamanByID)
	mux.HandleFunc("PUT /peminjaman/{id}", store.handlePutPeminjaman)
	mux.HandleFunc("DELETE /peminjaman/{id}", store.handleDeletePeminjaman)

	// Routing Persona E: Blokir
	mux.HandleFunc("POST /blokir", store.handlePostBlokir)
	mux.HandleFunc("GET /blokir", store.handleGetBlokir)
	mux.HandleFunc("DELETE /blokir/{id}", store.handleDeleteBlokir)

	return mux
}

func main() {
	handler := NewHandler()
	log.Println("Server berjalan di port 8080...")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}