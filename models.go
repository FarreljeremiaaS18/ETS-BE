package main

import "sync"

type SumberDaya struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Jenis string `json:"jenis"`
	Stok  int    `json:"stok"`
}

type Peminjaman struct {
	ID           int    `json:"id"`
	SumberDayaID int    `json:"sumber_daya_id"`
	Peminjam     string `json:"peminjam"`
	Tanggal      string `json:"tanggal"`
	SlotMulai    int    `json:"slot_mulai"`
	SlotSelesai  int    `json:"slot_selesai"`
	Jumlah       int    `json:"jumlah"`
	Status       string `json:"status"` // "aktif" atau "dibatalkan"
}

type Blokir struct {
	ID           int    `json:"id"`
	SumberDayaID int    `json:"sumber_daya_id"`
	Tanggal      string `json:"tanggal"`
	SlotMulai    int    `json:"slot_mulai"`
	SlotSelesai  int    `json:"slot_selesai"`
	Keterangan   string `json:"keterangan"`
}

type MemoryStore struct {
	mu sync.RWMutex
	SumberDaya map[int]SumberDaya
	Peminjaman map[int]Peminjaman
	Blokir     map[int]Blokir

	NextSumberDayaID int
	NextPeminjamanID int
	NextBlokirID     int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		SumberDaya:       make(map[int]SumberDaya),
		Peminjaman:       make(map[int]Peminjaman),
		Blokir:           make(map[int]Blokir),
		NextSumberDayaID: 1,
		NextPeminjamanID: 1,
		NextBlokirID:     1,
	}
}
