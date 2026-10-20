# PROSES

## 1. Persona dan masalahnya
(2-3 kalimat dengan kata-katamu sendiri: siapa pengguna persona kamu, masalah apa yang diselesaikan fitur khasnya.)
Pengguna persona saya adalah Asisten Laboratorium yang ingin mengamankan penggunaan lab saat praktikum maupun diservis. Masalah yang ditemukan adalah jika ada servis mendadak, lab tidak boleh dipinjam , dan jika sudah terlanjur dipinjam, peminjaman harus dibatalkan.

## 2. Rancangan struktur data
(Tulis struct dan slice/map apa yang kamu simpan, dan kenapa dipilih begitu. Rancang ini sebelum menulis kode, lalu catat kalau berubah di tengah jalan.)
Penyimpanan data dilakukan di memori menggunakan struktur MemoryStore yang di dalamnya memanfaatkan tipe data map[int]T untuk entitas SumberDaya, Peminjaman, dan Blokir. Pemilihan map dengan ID sebagai key digunakan untuk memastikan performa pencarian, pembaruan (PUT), dan penghapusan (DELETE) berjalan cepat dengan kompleksitas O(1) dibandingkan melakukan iterasi dengan slice. Setiap atribut dalam struct disesuaikan dengan tag JSON yang tepat berdasarkan kontrak API. Untuk memenuhi kebutuhan fitur konkurensi (B7), struktur MemoryStore membungkus seluruh data dengan sync.RWMutex, sehingga setiap operasi modifikasi (POST/PUT/DELETE) maupun validasi bentrok dapat dikunci dengan aman, mencegah dua peminjaman atau blokir bertabrakan di waktu yang bersamaan.


## 3. Satu keputusan desain
(Satu keputusan yang kamu ambil, pilihan lain yang kamu pertimbangkan, dan alasanmu memilih yang ini. Contoh: kenapa cek bentrok per slot, bukan per rentang.)
Untuk mendeteksi bentrok peminjaman, saya menghitung total penggunaan unit per slot menggunakan array berukuran 8 (mewakili slot 1-8) untuk suatu tanggal dan sumber daya tertentu.
Pilihan lain yang dipertimbangkan: Hanya melakukan pengecekan irisan (menggunakan rumus matematika range intersection) antar peminjaman satu per satu.
Alasan memilih ini: Karena sistem mengizinkan satu sumber daya dipinjam beberapa kali di slot yang sama asalkan total jumlahnya tidak melebihi stok, sekadar mengecek irisan rentang tidak cukup. Dengan mengisi array akumulasi pemakaian di tiap slot, saya bisa memastikan kapasitas slot tidak terlampaui dengan logika yang lebih terstruktur.

## 4. Bug atau kendala (minimal dua)
### Bug 1
- Gejala:
- Cara menemukan:
- Penyebab dan solusi:

### Bug 2
- Gejala:
- Cara menemukan:
- Penyebab dan solusi:

## 5. Change request
(Apa yang diminta, bagian kode mana yang kamu ubah, dan satu keputusan desain yang kamu ambil.)
Permintaan: Menambahkan field bentrok_dengan pada respons HTTP 409 khusus untuk POST /peminjaman.
Kode yang diubah: Memodifikasi fungsi cekBentrok di validasi.go untuk mengembalikan array ID, serta memodifikasi handlePostPeminjaman dan handlePutPeminjaman di handler.go.
Satu keputusan desain: Saya mengubah iterasi pengecekan dari range map menjadi for i := 1; i < s.NextPeminjamanID; i++. Karena map di Go tidak memiliki urutan, iterasi berbasis ID ini mengamankan syarat dokumen CR agar data bentrok_dengan terjamin urut naik dari awal tanpa perlu memanggil fungsi sort tambahan yang memakan proses komputasi.