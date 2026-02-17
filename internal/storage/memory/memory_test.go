package memory

import (
	"testing"
	"urlshort/internal/storage"	
)

func TestMemoryStorage_SaveAndGet(t *testing.T) {
	store := NewMemoryStorage()

	original := "https://example.com"
	short, err := store.Save(original)
	if err != nil{
		t.Fatalf("Save failed: %v", err)
	}

	short2, err := store.Save(original)
	if err != nil{
		t.Fatalf("Save duplicate failed: %v", err)
	}
	if short != short2{
		t.Errorf("expected same short for duplicate URL, got %s and %s", short, short2)
	}

	got, err := store.Get(short)
	if err != nil{
		t.Fatalf("Get failed: %v", err)
	}
	if got != original{
		t.Errorf("expected %s, got %s", original, got)
	}

	_, err = store.Get("not_exist")
	if err != storage.ErrNotFound{
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}