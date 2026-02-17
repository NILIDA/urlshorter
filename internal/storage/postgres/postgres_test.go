package postgres

import (
	"database/sql"
	"os"
	"testing"

	"urlshort/internal/storage"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func getTestConnStr() string {
	if dsn := os.Getenv("TEST_POSTGRES_DSN"); dsn != "" {
		return dsn
	}
	return "postgres://postgres:password@localhost:5432/shortener?sslmode=disable"
}

func clearTable(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM urls")
	return err
}

func mustConnect(tb testing.TB) *PostgresStorage {
	tb.Helper()
	connStr := getTestConnStr()
	store, err := NewPostgresStorage(connStr)
	if err != nil {
		tb.Skipf("can't connect to test database: %v", err)
	}
	return store
}

func TestPostgresStorage_SaveAndGet(t *testing.T) {
	store := mustConnect(t)
	defer store.Close()

	if err := clearTable(store.db); err != nil {
		t.Fatalf("failed to clear table: %v", err)
	}

	t.Run("save new url", func(t *testing.T) {
		original := "https://example.com"
		short, err := store.Save(original)
		if err != nil {
			t.Fatalf("Save failed: %v", err)
		}
		if len(short) != 10 {
			t.Errorf("expected short length 10, got %d", len(short))
		}

		got, err := store.Get(short)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if got != original {
			t.Errorf("expected %s, got %s", original, got)
		}
	})

	t.Run("save duplicate url returns same short", func(t *testing.T) {
		original := "https://duplicate.com"
		short1, err := store.Save(original)
		if err != nil {
			t.Fatalf("first Save failed: %v", err)
		}

		short2, err := store.Save(original)
		if err != nil {
			t.Fatalf("second Save failed: %v", err)
		}

		if short1 != short2 {
			t.Errorf("expected same short for duplicate url, got %s and %s", short1, short2)
		}
	})

	t.Run("get nonexistent short returns ErrNotFound", func(t *testing.T) {
		_, err := store.Get("nonexist")
		if err != storage.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestPostgresStorage_SaveConflictRetry(t *testing.T) {
	store := mustConnect(t)
	defer store.Close()

	if err := clearTable(store.db); err != nil {
		t.Fatalf("failed to clear table: %v", err)
	}

	_, err := store.Save("https://collision-test.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = store.Save("https://another-example.com")
	if err != nil {
		t.Errorf("unexpected error on second save: %v", err)
	}
}