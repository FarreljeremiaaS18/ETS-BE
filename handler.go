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

	if err := s.validateSumberDaya(&input); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
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

	if err := s.validatePeminjaman(&input); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	bentrokIDs, err := s.cekBentrok(&input, 0)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict) 
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":          err.Error(),
			"bentrok_dengan": bentrokIDs, 
		})
		return
	}

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
		sendError(w, http.StatusBadRequest, "ID bukan angka") 
		return
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	p, exists := s.Peminjaman[id]
	if !exists {
		sendError(w, http.StatusNotFound, "Peminjaman tidak ada") 
		return
	}

	sendJSON(w, http.StatusOK, p)
}

func (s *MemoryStore) handlePutPeminjaman(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "ID bukan angka")
		return
	}

	var input Peminjaman
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendError(w, http.StatusBadRequest, "Body tidak valid")
		return
	}
	
	if err := s.validatePeminjaman(&input); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	s.mu.Lock()
	defer s.mu.Unlock()
	
	oldData, exists := s.Peminjaman[id]
	if !exists {
		sendError(w, http.StatusNotFound, "Peminjaman tidak ada")
		return
	}
	
	if _, err := s.cekBentrok(&input, id); err != nil {
		sendError(w, http.StatusConflict, err.Error()) 
		return
	}
	
	input.ID = id
	input.Status = oldData.Status
	s.Peminjaman[id] = input

	
	sendJSON(w, http.StatusOK, input)
}

func (s *MemoryStore) handleDeletePeminjaman(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "ID bukan angka") 
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.Peminjaman[id]; !exists {
		sendError(w, http.StatusNotFound, "Peminjaman tidak ada")
		return
	}
	
	delete(s.Peminjaman, id)

	w.WriteHeader(http.StatusNoContent)
}


func (s *MemoryStore) handleGetBlokir(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := []Blokir{}
	
	for i := 1; i < s.NextBlokirID; i++ {
		if b, exists := s.Blokir[i]; exists {
			result = append(result, b)
		}
	}
	sendJSON(w, http.StatusOK, result)
}


func (s *MemoryStore) handleDeleteBlokir(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "ID bukan angka")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.Blokir[id]; !exists {
		sendError(w, http.StatusNotFound, "Blokir tidak ada")
		return
	}
	delete(s.Blokir, id)
	w.WriteHeader(http.StatusNoContent) 
}


func (s *MemoryStore) handlePostBlokir(w http.ResponseWriter, r *http.Request) {
	var input Blokir
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendError(w, http.StatusBadRequest, "Body tidak valid")
		return
	}

	if err := s.validateBlokir(&input); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()


	for _, b := range s.Blokir {
		if b.SumberDayaID == input.SumberDayaID && b.Tanggal == input.Tanggal {
			if input.SlotMulai <= b.SlotSelesai && input.SlotSelesai >= b.SlotMulai {
				sendError(w, http.StatusConflict, "beririsan dengan blokir lain") 
				return
			}
		}
	}

	
	input.ID = s.NextBlokirID
	s.Blokir[input.ID] = input
	s.NextBlokirID++

	
	dibatalkan := []int{} 
	for i := 1; i < s.NextPeminjamanID; i++ {
		if p, exists := s.Peminjaman[i]; exists {
			if p.Status == "aktif" && p.SumberDayaID == input.SumberDayaID && p.Tanggal == input.Tanggal {
				if p.SlotMulai <= input.SlotSelesai && p.SlotSelesai >= input.SlotMulai {
					p.Status = "dibatalkan"
					s.Peminjaman[i] = p 
					dibatalkan = append(dibatalkan, p.ID)
				}
			}
		}
	}

	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"blokir":     input,
		"dibatalkan": dibatalkan,
	})
}