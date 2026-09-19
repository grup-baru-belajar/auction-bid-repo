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
| USER | Login dan melakukan bid |

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
{
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
  "username": "john",
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
Content-Type: application/json
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
    "lastPrice": 10000000,
    "createdAt": "2026-09-18T08:00:00Z",
    "endTime": "2026-12-31T23:59:59Z",
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

```json
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

# 4. Get Auctions

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
Authorization: Bearer <token>
Content-Type: application/json
```

```json
{
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

# Database Schema

## Users

```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    username VARCHAR(100) NOT NULL UNIQUE,
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
| 403 | Forbidden |
| 404 | Resource Not Found |
| 409 | Conflict |
| 500 | Internal Server Error |