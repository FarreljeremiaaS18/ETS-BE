# AI-LOG

Catat setiap penggunaan AI yang menghasilkan kode. Jika tidak memakai AI sama sekali, tulis satu baris: **tidak memakai AI**.

Log yang jujur tidak mengurangi nilai, sebanyak apa pun AI dipakai. Yang mengurangi nilai adalah log yang tidak sesuai dengan kode, atau isi log yang tidak bisa dijelaskan.

## (Nama fitur atau fungsi, misalnya: Deteksi bentrok slot, B5)
- Prompt: Saya sedang mengerjakan ETS mata kuliah back-end, saya mendapat persona Asisten Laboratorium, bantu saya untuk membuat rancangan struktur data untuk project ini
- Dipakai: 
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
	Status       string `json:"status"`
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
	mu               sync.RWMutex
	SumberDaya       map[int]SumberDaya
	Peminjaman       map[int]Peminjaman
	Blokir           map[int]Blokir
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
- Diubah: -

## (Entri berikutnya)
- Prompt: Sekarang saya sudah selesai untuk mengerjakan setup untuk bagian struktur data, langkah selanjutnya setup untuk server dan routing. Jelaskan jangan hanya langsung beri kodenya tanpa penjelasan
- Dipakai: 
- Diubah:
