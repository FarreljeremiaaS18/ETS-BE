# AI-LOG

Catat setiap penggunaan AI yang menghasilkan kode. Jika tidak memakai AI sama sekali, tulis satu baris: **tidak memakai AI**.

Log yang jujur tidak mengurangi nilai, sebanyak apa pun AI dipakai. Yang mengurangi nilai adalah log yang tidak sesuai dengan kode, atau isi log yang tidak bisa dijelaskan.

## (Nama fitur atau fungsi, misalnya: Deteksi bentrok slot, B5)
# AI-LOG

## 1. Rancangan Struktur Data dan Setup Server (models.go & main.go)
- Prompt: Minta rancangan struktur data untuk Persona E (Asisten Laboratorium) dan cara setup server routing dasar.
- Dipakai: Struct `SumberDaya`, `Peminjaman`, `Blokir`, `MemoryStore` di `models.go` dan fungsi `NewHandler` di `main.go`.
- Diubah: Saya memecah kode hasil AI yang awalnya berpotensi digabung menjadi file-file terpisah (`models.go`, `main.go`, `handler.go`) agar lebih rapi dan mengatasi error "undefined method" saat *routing* dipanggil.

## 2. Deteksi bentrok slot dan Blokir Jadwal (validasi.go & handler.go)
- Prompt: Arahan untuk mengerjakan fitur deteksi bentrok (B5) dan fitur Blokir (Persona E).
- Dipakai: Fungsi `cekBentrok` dan `validateBlokir` di `validasi.go`, serta logika `handlePostBlokir` di `handler.go`.
- Diubah: Hasil AI yang menyarankan iterasi biasa saya integrasikan dengan keputusan desain saya (menggunakan *array* berukuran 9) untuk menghitung kapasitas per slot, serta mengatur agar pembatalan jadwal dieksekusi di dalam satu *lock* mutex yang sama.

## 3. Change Request B: Daftar Bentrok (validasi.go)
- Prompt: Arahan untuk mengerjakan lembar CR-B (Daftar Bentrok) yang mengembalikan array ID peminjaman yang menghalangi.
- Dipakai: Modifikasi logika pada fungsi `cekBentrok` untuk mengembalikan array `bentrokIDs`.
- Diubah: Saya menyempurnakan logika iterasi AI dengan mengubah iterasi `map` biasa menjadi perulangan `for i := 1; i < s.NextPeminjamanID; i++` agar ID `bentrok_dengan` yang dikembalikan sudah otomatis terurut naik tanpa perlu memanggil fungsi *sort* tambahan.

## 4. CRUD Dasar dan Validasi (Fitur B1-B4, B6)
- Prompt: Arahan langkah demi langkah untuk mengerjakan B1 (Sumber daya), B2-B3 (Peminjaman), B4 (Validasi error), dan B6 (Ubah dan hapus).
- Dipakai: Kerangka fungsi handler CRUD dasar di `handler.go` dan validasi input di `validasi.go`.
- Diubah: Menyesuaikan penamaan *field* JSON, pesan *error*, dan struktur respons agar mutlak sesuai dengan kontrak API di dokumen Requirement.