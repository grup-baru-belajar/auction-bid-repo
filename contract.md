<<<<<<< HEAD

=======
>>>>>>> bf56b32 (feat(backend): add bidding functionality with bid service and handler)
# Project Title

A brief description of what this project does and who it's for

# Auction Service API Contract

## Base URL

```http
/api/v1
```

---

# Authentication

Menggunakan JWT Bearer Token.

Request Header:

```http
Authorization: Bearer <jwt_token>
```

### Roles

| Role | Permission |
|--------|--------|
| ADMIN | Membuat auction |

---

# Standard Response Format

## Success Response

```json
{
  "success": true,
  "message": "Operation successful",
  "data": {}
}
```

## Error Response

```json
  "success": false,
  "message": "Error message"
}
```

---

# 1. Login

## POST /login

Login user dan mengembalikan JWT Token.

### Request

```http
POST /login
Content-Type: application/json
```

```json
{
  "password": "password123"
}
```

### Response 200

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "id": 1,
    "name": "John Doe",
    "username": "john",
    "role": "USER",
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

### Response 401

```json
{
  "success": false,
  "message": "Invalid username or password"
}
```

---

# 2. Create Auction

## POST /auctions

Membuat auction baru.

### Authorization

ADMIN only.

### Request

```http
POST /auctions
Authorization: Bearer <token>
```

```json
{
  "auctionName": "iPhone 15 Pro",
  "description": "Brand new iPhone 15 Pro 256GB Natural Titanium",
  "imageLink": "https://cdn.example.com/images/iphone15pro.jpg",
  "startingPrice": 10000000,
  "endTime": "2026-12-31T23:59:59Z"
}

```

### Validation Rules

- auctionName wajib diisi
- startingPrice > 0
- endTime harus lebih besar dari waktu saat ini

### Response 201

```json
{
  "success": true,
  "message": "Auction created successfully",
  "data": {
    "id": 1,
    "auctionName": "iPhone 15 Pro",
    "description": "Brand new iPhone 15 Pro 256GB Natural Titanium",
    "imageLink": "https://cdn.example.com/images/iphone15pro.jpg",
    "startingPrice": 10000000,
    "isCompleted": false,
    "bidWinner": null
  }
}
```

### Response 403

```json
{
  "success": false,
  "message": "Only admin can create auction"
}
```

---

# 3. Get Auction Detail

## GET /auctions/{id}

Mengambil detail auction berdasarkan ID.

### Request

```http
GET /auctions/1
```

### Response 200

{
  "success": true,
  "message": "Auction retrieved successfully",
  "data": {
    "id": 1,
    "auctionName": "iPhone 15 Pro",
    "description": "Brand new iPhone 15 Pro 256GB Natural Titanium",
    "imageLink": "https://cdn.example.com/images/iphone15pro.jpg",
    "startingPrice": 10000000,
    "lastPrice": 15000000,
    "createdAt": "2026-09-18T08:00:00Z",
    "endTime": "2026-12-31T23:59:59Z",
    "isCompleted": false,
    "bidWinner": {
      "id": 5,
      "name": "John Doe"
    },
<<<<<<< HEAD
    "totalBids": 12,
    "totalBidders": 5,
=======
>>>>>>> bf56b32 (feat(backend): add bidding functionality with bid service and handler)
    "topBids": [
      {
        "id": 10,
        "userId": 5,
        "userName": "John Doe",
        "bidPrice": 15000000,
        "createdAt": "2026-09-18T08:20:00Z"
      },
      {
        "id": 8,
        "userId": 3,
        "userName": "Jane Smith",
        "bidPrice": 14500000,
        "createdAt": "2026-09-18T08:15:00Z"
      },
      {
        "id": 7,
        "userId": 2,
        "userName": "Mike Lee",
        "bidPrice": 14000000,
        "createdAt": "2026-09-18T08:10:00Z"
      }
    ]
  }
}
```

### Notes

- `topBids` diurutkan berdasarkan `bidPrice DESC`.
- Maksimum hanya menampilkan 3 bid tertinggi.

### Response 404

```json
{
  "success": false,
  "message": "Auction not found"
}
```

---

## GET /auctions

Mengambil daftar seluruh auction.

### Request

```http
GET /auctions
```
### Query Parameters

| Parameter | Type | Required | Description |
|------------|------------|------------|------------|
| page | integer | No | Default 1 |
| limit | integer | No | Default 10 |
| isCompleted | boolean | No | Filter status auction |

### Sample Request

```http
GET /auctions?page=1&limit=10
```

### Response 200

```json
{
  "success": true,
  "message": "Auctions retrieved successfully",
  "data": [
    {
      "id": 1,
      "auctionName": "iPhone 15 Pro",
      "description": "Brand new iPhone 15 Pro 256GB Natural Titanium",
      "imageLink": "https://cdn.example.com/images/iphone15pro.jpg",
      "startingPrice": 10000000,
      "lastPrice": 15000000,
      "createdAt": "2026-09-18T08:00:00Z",
      "endTime": "2026-12-31T23:59:59Z",
      "isCompleted": false,
      "bidWinner": {
        "id": 5,
        "name": "John Doe"
      }
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 1,
    "totalPages": 1
  }
}
```

---

# 5. Create Bid

## POST /bid

Membuat bid baru pada auction.

### Authorization

Login required.

### Business Rules
- User harus login.
- Auction harus masih aktif.
- Auction belum selesai (`isCompleted = false`).
- Waktu sekarang belum melewati `endTime`.
- BidPrice harus lebih besar dari `lastPrice`.
- Saat bid berhasil:
  - `Auctions.lastPrice` diperbarui.
  - `Auctions.bidWinnerID` diperbarui dengan UserID bidder terbaru.

### Request

```http
POST /bid
  "auctionId": 1,
  "bidPrice": 16000000
}
```

### Response 201

```json
{
  "success": true,
  "message": "Bid placed successfully",
  "data": {
    "id": 15,
    "auctionId": 1,
    "userId": 5,
    "bidPrice": 16000000,
    "createdAt": "2026-09-18T08:20:00Z"
  }
}
```

### Response 400

```json
{
  "success": false,
  "message": "Bid price must be greater than current price"
}
```

### Response 401

```json
{
  "success": false,
  "message": "Unauthorized"
}
```

### Response 404

```json
{
  "success": false,
  "message": "Auction not found"
}
```

### Response 409

```json
{
  "success": false,
  "message": "Auction already completed"
}
```

---

<<<<<<< HEAD
# 6. Reporting

Endpoint untuk dashboard reporting/analytics. Semua endpoint di bawah ini **ADMIN only** dan memerlukan JWT Bearer Token.

## 6.1 Get Top Auctions

### GET /reporting/top-auction

Mengambil 5 auction dengan jumlah bid terbanyak.


ADMIN only.

### Request

```http
GET /reporting/top-auction
Authorization: Bearer <token>
```

### Response 200

```json
{
  "success": true,
  "message": "Top auctions retrieved successfully",
  "data": [
    {
      "id": 1,
      "auctionName": "iPhone 15 Pro",
      "startingPrice": 10000000,
      "highestBid": 16000000,
      "totalBid": 12,
      "bidders": 5,
      "status": "ACTIVE"
    }
  ]
}
```

### Notes

- `status` bernilai `"ACTIVE"` atau `"ENDED"`, dihitung dari `is_completed` dan `end_time`.
- Diurutkan berdasarkan `totalBid DESC`, dibatasi 5 hasil.

---

## 6.2 Get Auction Activity

### GET /reporting/auction-activity

Mengambil jumlah bid per hari untuk 7 hari terakhir.

### Authorization


```http
GET /reporting/auction-activity
Authorization: Bearer <token>
```

### Response 200

```json
{
  "success": true,
  "message": "Auction activity retrieved successfully",
  "data": [
    {
      "date": "2026-09-14",
      "totalBids": 8
    },
    {
      "date": "2026-09-15",
      "totalBids": 3
    }
  ]
}
```

### Notes

- Hanya mengembalikan tanggal yang memiliki bid (tidak mengisi tanggal kosong dengan 0).
- Diurutkan `date ASC`.

---

## 6.3 Get Auction Status

### GET /reporting/auction-status

Mengambil jumlah auction per status (ACTIVE / ENDED).

### Authorization


```http
GET /reporting/auction-status
Authorization: Bearer <token>
```

### Response 200

```json
{
  "success": true,
  "message": "Auction status retrieved successfully",
  "data": [
    {
      "status": "ACTIVE",
      "total": 7
    },
    {
      "status": "ENDED",
      "total": 3
    }
  ]
}
```

---

## 6.4 Get Auction Summary

### GET /reporting/auction-summary

Mengambil ringkasan agregat seluruh auction dan bid.

### Authorization

ADMIN only.


```http
GET /reporting/auction-summary
Authorization: Bearer <token>
```

### Response 200

```json
{
  "success": true,
  "message": "Auction summary retrieved successfully",
  "data": {
    "totalAuctions": 10,
    "ongoingAuctions": 7,
    "completedAuctions": 3,
    "totalBidsOngoing": 25,
    "totalBidsCompleted": 15,
    "totalBidsAll": 40
  }
}
```

---

## Reporting Error Responses

### Response 401

```json
{
  "success": false,
  "message": "Unauthorized"
}
```
### Response 403

```json
{
  "success": false,
  "message": "Only admin can create auction"
}
```

### Notes

- Pesan 403 di atas memang generik (dipakai ulang dari middleware `RequireAdmin`), bukan pesan khusus reporting.
- Response 500 mengikuti format standar dengan pesan spesifik per endpoint, contoh: `"Failed to get auction summary"`.

---

=======
>>>>>>> bf56b32 (feat(backend): add bidding functionality with bid service and handler)
# Database Schema

## Users

```sql
    name VARCHAR(100) NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL
);
```

## Auctions

```sql
CREATE TABLE auctions (
    id BIGSERIAL PRIMARY KEY,
    auction_name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    image_link TEXT NOT NULL,
    bid_winner_id BIGINT NULL,
    starting_price NUMERIC(18,2) NOT NULL,
    last_price NUMERIC(18,2) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    is_completed BOOLEAN DEFAULT FALSE,

    CONSTRAINT fk_auction_bid_winner
        FOREIGN KEY (bid_winner_id)
        REFERENCES users(id)
);
```

## Bid

```sql
CREATE TABLE bid (
    id BIGSERIAL PRIMARY KEY,
    auction_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    bid_price NUMERIC(18,2) NOT NULL,
    created_at TIMESTAMP NOT NULL,

    CONSTRAINT fk_bid_auction
        FOREIGN KEY (auction_id)
        REFERENCES auctions(id),

    CONSTRAINT fk_bid_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
);
```

---

# Error Codes

| HTTP Status | Description |
|------------|------------|
| 200 | Success |
| 201 | Created |
| 400 | Validation Error |
| 401 | Unauthorized |
| 404 | Resource Not Found |
| 409 | Conflict |
| 500 | Internal Server Error |