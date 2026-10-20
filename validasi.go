package main

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)


func (s *MemoryStore) validateSumberDaya(input *SumberDaya) error {
	
	input.Nama = strings.TrimSpace(input.Nama)
	
	if input.Nama == "" {
		return errors.New("nama tidak boleh kosong")
	}
	
	if utf8.RuneCountInString(input.Nama) > 50 {
		return errors.New("nama maksimal 50 karakter")
	}
	if input.Jenis != "ruang" && input.Jenis != "alat" {
		return errors.New("jenis hanya boleh ruang atau alat")
	}
	if input.Stok < 1 {
		return errors.New("stok minimal 1")
	}
	return nil
}


func (s *MemoryStore) validatePeminjaman(input *Peminjaman) error {
	
	input.Peminjam = strings.TrimSpace(input.Peminjam)
	
	if input.Peminjam == "" {
		return errors.New("peminjam tidak boleh kosong")
	}
	if utf8.RuneCountInString(input.Peminjam) > 50 {
		return errors.New("peminjam maksimal 50 karakter")
	}

	
	if _, err := time.Parse("2006-01-02", input.Tanggal); err != nil {
		return errors.New("format tanggal salah atau tanggal tidak ada")
	}

	
	if input.SlotMulai < 1 || input.SlotSelesai > 8 || input.SlotMulai > input.SlotSelesai {
		return errors.New("rentang slot tidak valid")
	}

	
	s.mu.RLock()
	sd, exists := s.SumberDaya[input.SumberDayaID]
	s.mu.RUnlock()

	if !exists {
		
		return errors.New("sumber daya tidak ditemukan") 
	}
	if input.Jumlah < 1 || input.Jumlah > sd.Stok {
		return errors.New("jumlah harus antara 1 dan stok sumber daya")
	}

	return nil
}

var ErrBentrok = errors.New("bentrok kapasitas slot")


func (s *MemoryStore) cekBentrok(input *Peminjaman, ignoreID int) error {
	
	sd := s.SumberDaya[input.SumberDayaID]

	for _, b := range s.Blokir {
		if b.SumberDayaID == input.SumberDayaID && b.Tanggal == input.Tanggal {
			if input.SlotMulai <= b.SlotSelesai && input.SlotSelesai >= b.SlotMulai {
				return errors.New("jadwal diblokir oleh asisten laboratorium") // 409[cite: 1]
			}
		}
	}

	terpakai := make([]int, 9)

	
	for _, p := range s.Peminjaman {
		
		if p.Status == "aktif" && p.SumberDayaID == input.SumberDayaID && p.Tanggal == input.Tanggal {
			
			if p.ID == ignoreID {
				continue 
			}
			
			for i := p.SlotMulai; i <= p.SlotSelesai; i++ {
				terpakai[i] += p.Jumlah
			}
		}
	}

	
	for i := input.SlotMulai; i <= input.SlotSelesai; i++ {
		if terpakai[i]+input.Jumlah > sd.Stok {
			return ErrBentrok 
		}
	}

	return nil
}

func (s *MemoryStore) validateBlokir(input *Blokir) error {
	input.Keterangan = strings.TrimSpace(input.Keterangan)
	if input.Keterangan == "" {
		return errors.New("keterangan tidak boleh kosong")
	}
	if utf8.RuneCountInString(input.Keterangan) > 100 {
		return errors.New("keterangan maksimal 100 karakter")
	}
	if _, err := time.Parse("2006-01-02", input.Tanggal); err != nil {
		return errors.New("tanggal tidak valid")
	}
	if input.SlotMulai < 1 || input.SlotSelesai > 8 || input.SlotMulai > input.SlotSelesai {
		return errors.New("rentang slot tidak valid")
	}
	
	s.mu.RLock()
	_, exists := s.SumberDaya[input.SumberDayaID]
	s.mu.RUnlock()
	if !exists {
		return errors.New("sumber daya tidak ditemukan") // 400
	}
	return nil
}