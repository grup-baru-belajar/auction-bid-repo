# Bruno Collection — Auction API

Collection untuk menguji endpoint **login** dan **get auction by id** lewat HTTP,
mencakup kasus positive, negative, dan edge yang sama dengan unit test Go-nya.

## Cara pakai

### Lewat aplikasi Bruno
1. Buka Bruno → **Open Collection** → pilih folder `bruno/` ini
2. Pilih environment **Local** di pojok kanan atas
3. Jalankan folder `01-login` dulu (request pertama menyimpan token untuk folder 03)

### Lewat terminal
```bash
cd bruno
npx @usebruno/cli run --env Local -r
```

## Prasyarat

Backend harus hidup dan database harus sudah di-seed:

```bash
cd ../backend
docker compose up -d postgres
go run . serve
```

Seed menyediakan `admin` / `admin123` dan `john` / `password123`.

## Isi

| Folder | Jumlah | Cakupan |
|---|---|---|
| `01-login` | 10 | happy path admin & user, kredensial salah, validasi body, JSON rusak |
| `02-auction-detail` | 7 | detail lengkap, struktur topBids, 404, dan 4 bentuk id tidak valid |
| `03-auth-middleware` | 6 | tanpa token, skema salah, token ngawur, signature dipalsukan, 403 vs 200 |

Urutannya penting: `01-login/01` menyimpan `{{adminToken}}`, `01-login/02` menyimpan
`{{userToken}}`, dan `01-login/03` menyimpan pesan error yang dibandingkan oleh `01-login/04`.

## Catatan

Tes Go punya satu kasus yang tidak bisa ditiru di sini: body terpotong `{"username":`.
Kurung kurawalnya tidak seimbang, dan parser `.bru` ikut menelan kurung penutup bloknya
sehingga body berubah jadi `{"username":}` yang memicu cabang error berbeda.
Cabang yang sama sudah ditutup oleh `10-body-kosong`, yang juga menghasilkan
`Invalid request body`.
