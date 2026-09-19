# Backend

API Auction

## Yang perlu disiapkan

- Go 1.27+
- Docker (untuk database)

## 1. Siapkan file config

Dua file ini tidak ikut di Git, jadi bikin sendiri:

```bash
cp .env.example .env
cp config.example.yaml config.yaml
```

Isi `.env` sesuai selera. `JWT_SECRET` wajib minimal 32 karakter — kalau mau
bikin yang acak:

```bash
openssl rand -base64 48 | tr -d '/+=' | head -c 48
```

## 2. Nyalakan database

```bash
docker compose up -d postgres
```

Cek sudah siap:

```bash
docker compose ps
```

## 3. Jalankan migration

```bash
docker compose exec -T postgres \
  sh -c 'psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$POSTGRES_DB"' \
  < migrations/000001_init_schema.up.sql
```


Untuk membatalkan (hapus semua tabel), ganti `up.sql` jadi `down.sql`.

Ini hanya membuat tabel di database

## 4. Isi data contoh

```bash
docker compose exec -T postgres \
  sh -c 'psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$POSTGRES_DB"' \
  < seed/dev_seed.sql
```

Isinya 5 user, 5 lelang, dan 11 bid. Passwordnya `admin` → `admin123`,
sisanya → `password123`.


## 5. Jalankan servernya

```bash
go run . serve
```

Kalau berhasil:

```
app.env=development app.port=8080
database: connected
http: listening on :8080
```

Coba login:

```bash
curl -X POST http://localhost:8080/api/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"admin123"}'
```

Hentikan dengan Ctrl+C. Server menunggu request yang sedang jalan selesai dulu.

## Adminer (opsional)

Buat lihat isi database lewat browser:

```bash
docker compose up -d adminer
```

Buka http://localhost:8081, lalu isi:

| Field | Isi |
|---|---|
| System | PostgreSQL |
| Server | `postgres` |
| Username | dari `.env` |
| Password | dari `.env` |
| Database | dari `.env` |

