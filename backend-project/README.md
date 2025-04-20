# Product Management Backend (Golang + Gin)

Project ini merupakan backend sederhana untuk manajemen produk, termasuk fitur upload, download, dan delete gambar produk.

## Struktur Folder
```
.
├── config/
│   └── db.go
├── handlers/
│   └── product_handler.go
├── routes/
│   └── router.go
├── uploads/
├── main.go
```
- **config/**: Konfigurasi koneksi database MySQL
- **handlers/**: Handler untuk upload, download, delete gambar produk
- **routes/**: Routing API
- **uploads/**: Folder penyimpanan file gambar
- **main.go**: Entry point aplikasi

## Fitur
- Upload gambar produk (.png, .jpg, .jpeg)
- Download gambar produk berdasarkan ID
- Delete gambar produk
- Validasi ukuran file maksimum 5MB
- Validasi format file gambar
- Error handling lengkap

## Endpoint API
| Endpoint | Method | Deskripsi |
| :--- | :--- | :--- |
| `/products/:id/upload` | POST | Upload gambar produk |
| `/products/:id/image` | GET | Download gambar produk |
| `/products/:id/image` | DELETE | Hapus gambar produk |

## Cara Menjalankan
1. Clone project ini.
2. Pastikan sudah install Golang, MySQL, dan Gin framework (`go get github.com/gin-gonic/gin`).
3. Setup database MySQL dan buat tabel `products` minimal dengan field `id` (INT) dan `image_url` (VARCHAR).
4. Atur koneksi database di `config/db.go`.
5. Jalankan server:
    ```bash
    go run main.go
    ```
6. Test API menggunakan Postman (Collection sudah disediakan).

## Catatan
- File yang diupload akan disimpan di folder `uploads/`.
- Pastikan folder `uploads/` memiliki permission tulis.
- Ukuran file dibatasi maksimal 5MB.

## Tools
- Golang 1.18+
- Gin Web Framework
- MySQL
- Postman (untuk testing)

---