package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)


func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json") 
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}


func sendError(w http.ResponseWriter, status int, message string) {
	sendJSON(w, status, map[string]string{"error": message}) 
}

func (s *MemoryStore) handlePostSumberDaya(w http.ResponseWriter, r *http.Request) {
	var input SumberDaya
	
	
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendError(w, http.StatusBadRequest, "Body tidak valid") 
		return
	}

	
	s.mu.Lock()
	defer s.mu.Unlock()


	input.ID = s.NextSumberDayaID
	s.SumberDaya[input.ID] = input
	s.NextSumberDayaID++

	
	sendJSON(w, http.StatusCreated, input)
}

func (s *MemoryStore) handleGetSumberDaya(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock() 
	defer s.mu.RUnlock()

	result := []SumberDaya{} 
	
	for _, sd := range s.SumberDaya {
		result = append(result, sd)
	}

	sendJSON(w, http.StatusOK, result)
}

func (s *MemoryStore) handleGetSumberDayaByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	
	if err != nil {
		sendError(w, http.StatusBadRequest, "id bukan angka") 
		return
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	
	sd, exists := s.SumberDaya[id]
	if !exists {
		sendError(w, http.StatusNotFound, "Sumber daya tidak ditemukan") 
		return
	}

	
	sendJSON(w, http.StatusOK, sd)
}


func (s *MemoryStore) handlePostPeminjaman(w http.ResponseWriter, r *http.Request) {
	var input Peminjaman
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendError(w, http.StatusBadRequest, "Body tidak valid") 
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	input.ID = s.NextPeminjamanID
	input.Status = "aktif" 
	s.Peminjaman[input.ID] = input
	s.NextPeminjamanID++

	sendJSON(w, http.StatusCreated, input) 
}


func (s *MemoryStore) handleGetPeminjaman(w http.ResponseWriter, r *http.Request) {
	
	filterSumberDaya := r.URL.Query().Get("sumber_daya_id")
	filterTanggal := r.URL.Query().Get("tanggal")
	filterStatus := r.URL.Query().Get("status")

	
	if filterSumberDaya != "" {
		if _, err := strconv.Atoi(filterSumberDaya); err != nil {
			sendError(w, http.StatusBadRequest, "sumber_daya_id bukan angka") 
			return
		}
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	result := []Peminjaman{} 
	
	for _, p := range s.Peminjaman {
		match := true
		
		
		if filterSumberDaya != "" && strconv.Itoa(p.SumberDayaID) != filterSumberDaya {
			match = false
		}
		if filterTanggal != "" && p.Tanggal != filterTanggal {
			match = false
		}
		if filterStatus != "" && p.Status != filterStatus {
			match = false
		}
		
		if match {
			result = append(result, p)
		}
	}

	sendJSON(w, http.StatusOK, result)
}

func (s *MemoryStore) handleGetPeminjamanByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "ID bukan angka") // Urutan validasi pertama[cite: 2]
		return
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	p, exists := s.Peminjaman[id]
	if !exists {
		sendError(w, http.StatusNotFound, "Peminjaman tidak ada") // 404 jika tidak ketemu[cite: 2]
		return
	}

	sendJSON(w, http.StatusOK, p)
}