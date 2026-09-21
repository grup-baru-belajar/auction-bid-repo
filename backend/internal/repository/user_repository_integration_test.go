//go:build integration
// +build integration

package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/grup-baru-belajar/auction-bid-repo/internal/models"
	_ "github.com/lib/pq"
)

func TestUserRepository_FindByUsername_EndToEnd(t *testing.T) {
	dsn := os.Getenv("PG_DSN")
	if dsn == "" {
		t.Skip("PG_DSN not set; skipping integration test")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	// Use single connection so SET search_path applies for all statements
	db.SetMaxOpenConns(1)

	schema := fmt.Sprintf("test_schema_%d", time.Now().UnixNano())
	if _, err := db.ExecContext(context.Background(), fmt.Sprintf("CREATE SCHEMA %s", schema)); err != nil {
		t.Fatalf("create schema: %v", err)
	}
	// ensure cleanup
	defer func() {
		_, _ = db.ExecContext(context.Background(), fmt.Sprintf("DROP SCHEMA %s CASCADE", schema))
	}()

	if _, err := db.ExecContext(context.Background(), fmt.Sprintf("SET search_path TO %s", schema)); err != nil {
		t.Fatalf("set search_path: %v", err)
	}

	// create minimal tables required by user repository
	createSQL := `
    CREATE TABLE users (
        id BIGSERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        username VARCHAR(100) NOT NULL UNIQUE,
        password VARCHAR(255) NOT NULL,
        role VARCHAR(20) NOT NULL
    );
    `

	if _, err := db.ExecContext(context.Background(), createSQL); err != nil {
		t.Fatalf("create tables: %v", err)
	}

	// insert test user
	var userID int64
	err = db.QueryRowContext(context.Background(),
		`INSERT INTO users (name, username, password, role) VALUES ($1,$2,$3,$4) RETURNING id`,
		"John Doe", "john", "$2a$04$hashed", string(models.RoleUser)).Scan(&userID)
	if err != nil {
		t.Fatalf("insert user: %v", err)
	}

	repo := NewUserRepository(db)

	t.Run("found", func(t *testing.T) {
		user, err := repo.FindByUsername(context.Background(), "john")
		if err != nil {
			t.Fatalf("find by username: %v", err)
		}
		if user.ID != userID || user.Name != "John Doe" || user.Username != "john" {
			t.Fatalf("unexpected user returned: %+v", user)
		}
		if user.Password != "$2a$04$hashed" {
			t.Fatalf("expected password hash to be loaded, got %q", user.Password)
		}
		if user.Role != models.RoleUser {
			t.Fatalf("expected role %q, got %q", models.RoleUser, user.Role)
		}
	})

	t.Run("not found", func(t *testing.T) {
		user, err := repo.FindByUsername(context.Background(), "ghost")
		if user != nil {
			t.Fatalf("expected nil user, got %+v", user)
		}
		if !errors.Is(err, ErrUserNotFound) {
			t.Fatalf("expected error %v, got %v", ErrUserNotFound, err)
		}
	})

	t.Run("username is case sensitive", func(t *testing.T) {
		_, err := repo.FindByUsername(context.Background(), "JOHN")
		if !errors.Is(err, ErrUserNotFound) {
			t.Fatalf("expected error %v, got %v", ErrUserNotFound, err)
		}
	})
}
