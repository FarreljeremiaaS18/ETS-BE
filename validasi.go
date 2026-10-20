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


func (s *MemoryStore) cekBentrok(input *Peminjaman, ignoreID int) ([]int, error) {
	for i := 1; i < s.NextBlokirID; i++ {
		if b, exists := s.Blokir[i]; exists {
			if b.SumberDayaID == input.SumberDayaID && b.Tanggal == input.Tanggal {
				if input.SlotMulai <= b.SlotSelesai && input.SlotSelesai >= b.SlotMulai {
					return []int{}, errors.New("jadwal diblokir oleh asisten laboratorium")
				}
			}
		}
	}

	sd := s.SumberDaya[input.SumberDayaID]
	terpakai := make([]int, 9)
	bentrokIDs := []int{} 

	for i := 1; i < s.NextPeminjamanID; i++ {
		if p, exists := s.Peminjaman[i]; exists {
			if p.Status == "aktif" && p.SumberDayaID == input.SumberDayaID && p.Tanggal == input.Tanggal {
				if p.ID == ignoreID {
					continue
				}
				
				isOverlap := false
				for j := p.SlotMulai; j <= p.SlotSelesai; j++ {
					terpakai[j] += p.Jumlah
					if j >= input.SlotMulai && j <= input.SlotSelesai {
						isOverlap = true
					}
				}
				if isOverlap {
					bentrokIDs = append(bentrokIDs, p.ID)
				}
			}
		}
	}

	for i := input.SlotMulai; i <= input.SlotSelesai; i++ {
		if terpakai[i]+input.Jumlah > sd.Stok {
			return bentrokIDs, ErrBentrok
		}
	}

	return nil, nil
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