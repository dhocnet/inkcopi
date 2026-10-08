# Inkcopi

**Inkcopi (Inkscape Color Picker)** adalah color picker berbasis TUI yang dibuat menggunakan Go.

Inkcopi dirancang untuk membantu menentukan warna ketika membuat aplikasi berbasis terminal/TUI. Karena color picker ini sendiri berjalan di terminal, warna dapat dilihat langsung dalam lingkungan yang nantinya digunakan oleh aplikasi TUI.

Inkcopi dapat membaca palet warna Inkscape dan menyediakan mode untuk mencampur warna secara manual menggunakan RGB maupun CMYK.

## Fitur

- Membaca palet warna dari Inkscape.
- Berganti-ganti palet Inkscape menggunakan tombol `P`.
- Menampilkan warna sebagai blok warna di terminal.
- Navigasi warna menggunakan Arrow Keys atau `HJKL`.
- Menampilkan informasi warna:
  - HEX
  - RGB
  - CMYK
- Manual RGB color mixer.
- Manual CMYK color mixer.
- Preview warna secara realtime.
- Berganti mode menggunakan tombol `Tab`.
- Berjalan sepenuhnya melalui terminal/TUI.
- Mendukung perubahan ukuran terminal secara dinamis.

## Mode

Inkcopi memiliki tiga mode utama:

1. **Inkscape Palette**
2. **Manual RGB**
3. **Manual CMYK**

Gunakan tombol `Tab` untuk berpindah antar-mode.

### Inkscape Palette

Mode ini membaca file palette `.gpl` yang tersedia pada instalasi Inkscape.

Gunakan:

- `Arrow Keys` atau `HJKL` — memilih warna
- `P` — berganti palet

Warna yang dipilih akan ditampilkan bersama informasi HEX, RGB, dan CMYK.

### Manual RGB

Mode ini digunakan untuk membuat warna RGB secara manual.

- `Up/Down` atau `K/J` — memilih channel
- `Left/Right` atau `H/L` — mengubah nilai sebesar 1
- `Shift+Left/Right` atau `Shift+H/L` — mengubah nilai sebesar 10

Channel yang tersedia:

- Red
- Green
- Blue

### Manual CMYK

Mode ini digunakan untuk membuat warna CMYK secara manual.

- `Up/Down` atau `K/J` — memilih channel
- `Left/Right` atau `H/L` — mengubah nilai sebesar 1
- `Shift+Left/Right` atau `Shift+H/L` — mengubah nilai sebesar 10

Channel yang tersedia:

- Cyan
- Magenta
- Yellow
- Black

## Persyaratan

Inkcopi ditujukan untuk lingkungan Linux.

Yang dibutuhkan:

- Linux
- Go 1.26 atau lebih baru
- Inkscape

Inkscape diperlukan apabila ingin menggunakan palet warna Inkscape. Jika palet Inkscape tidak ditemukan, Inkcopi akan menggunakan fallback palette sederhana.

## Menyiapkan Go

Pastikan Go sudah terpasang pada sistem.

Periksa versi Go dengan:
```

go version

```

Contoh:
```

go version go1.26.4 linux/amd64

```

Jika Go belum terpasang, silakan instal Go terlebih dahulu sesuai distribusi Linux yang digunakan.

Setelah Go tersedia, clone repository Inkcopi:
```

git clone https://github.com/dhocnet/inkcopi.git cd inkcopi

```

Kemudian download dependency:
```

go mod download

```

## Menjalankan dengan `go run`

Inkcopi dapat langsung dijalankan tanpa melakukan proses compile terlebih dahulu.

Dari direktori repository, jalankan:
```

go run .

```

Go akan mengompilasi program secara sementara kemudian menjalankannya.

Cara ini cocok jika hanya ingin mencoba Inkcopi atau sedang melakukan pengembangan.

Untuk keluar dari Inkcopi, tekan:
```

Q

```

atau:
```

Ctrl+C

```

## Compile sebagai Standalone Binary

Jika ingin menggunakan Inkcopi tanpa menjalankan `go run` setiap kali, program dapat dikompilasi menjadi standalone binary.

Gunakan:
```

go build -o inkcopi .

```

Setelah proses selesai, akan terdapat binary:
```

inkcopi

```

Jalankan dengan:
```

./inkcopi

```

Binary tersebut dapat dijalankan tanpa membutuhkan source code dan tanpa menjalankan `go run`.

### Instalasi ke PATH

Jika ingin menjalankan Inkcopi dari direktori mana pun, binary dapat dipindahkan ke salah satu direktori yang terdapat pada `PATH`, misalnya:
```

sudo install -m 755 inkcopi /usr/local/bin/inkcopi

```

Setelah itu cukup jalankan:
```

inkcopi

```

## Struktur Dependency

Inkcopi menggunakan beberapa library dari ekosistem Charm.

Library utama yang digunakan:

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — framework TUI.
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — styling terminal.

Dependency proyek sudah didefinisikan di `go.mod`, sehingga tidak perlu menginstalnya secara manual.

## Palet Inkscape

Inkcopi membaca file palette dengan format `.gpl`.

Pada Linux, Inkcopi terlebih dahulu mencoba mencari palette pada:
```

/usr/share/inkscape/palettes/

```

Jika tidak ditemukan, Inkcopi juga mencoba:
```

./palettes/

```

Karena itu, untuk mendapatkan palette Inkscape secara lengkap, pastikan Inkscape sudah terpasang pada sistem.

## Screenshot

### Inkscape Palette

Inkcopi menampilkan palette Inkscape sebagai kumpulan blok warna yang dapat dinavigasi menggunakan keyboard.

### Manual RGB

Mode RGB memungkinkan nilai Red, Green, dan Blue diubah secara langsung dengan preview warna yang diperbarui secara realtime.

### Manual CMYK

Mode CMYK memungkinkan nilai Cyan, Magenta, Yellow, dan Black diubah secara langsung dengan preview warna.

## Keyboard Shortcuts

| Tombol | Fungsi |
| --- | --- |
| `Tab` | Berganti mode |
| `P` | Berganti palet Inkscape |
| `Arrow Keys` | Navigasi / mengubah nilai |
| `HJKL` | Navigasi / mengubah nilai |
| `Shift + Left/Right` | Mengubah nilai sebesar 10 |
| `Q` | Keluar |
| `Ctrl+C` | Keluar |

## Tujuan

Inkcopi bukan ditujukan untuk menggantikan color picker pada aplikasi grafis.

Tujuan utamanya sederhana: membantu menentukan warna untuk aplikasi TUI dengan melihat warna tersebut langsung dari terminal.

Daripada memilih warna melalui aplikasi grafis kemudian memperkirakan bagaimana warna tersebut akan terlihat di terminal, Inkcopi memungkinkan proses tersebut dilakukan langsung di lingkungan terminal.

## Dibuat Menggunakan Go dan AI

Inkcopi ditulis menggunakan Go dan dikembangkan dengan bantuan Gemini AI.

Sebagian besar kode program dibuat dengan bantuan AI berdasarkan requirement dan perilaku aplikasi yang ditentukan selama proses pengembangan.

