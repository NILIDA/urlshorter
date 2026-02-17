package handlers

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"urlshort/internal/storage"
	"urlshort/internal/storage/memory"
)

const testBaseURL = "http://localhost:8080"

func newTestHandler(store storage.Storage, logger *zap.SugaredLogger, baseURL string) *Handler {
	return &Handler{
		store:   store,
		logger:  logger,
		baseURL: baseURL,
		tmpl:    template.New(""),
	}
}

func TestHandlerSave(t *testing.T) {
	store := memory.NewMemoryStorage()
	logger := zap.NewNop().Sugar()
	h := newTestHandler(store, logger, testBaseURL)

	reqBody := SaveRequest{URL: "https://example.com"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Save(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var resp SaveResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp.Short) != 10 {
		t.Errorf("expected short code length 10, got %d", len(resp.Short))
	}
}

func TestHandlerGet(t *testing.T) {
	store := memory.NewMemoryStorage()
	logger := zap.NewNop().Sugar()
	h := newTestHandler(store, logger, testBaseURL)

	originalURL := "https://example.com"
	short, err := store.Save(originalURL)
	if err != nil {
		t.Fatalf("failed to save test url: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/"+short, nil)
	req.Header.Set("Accept", "application/json")
	w := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/{short}", h.Get).Methods(http.MethodGet)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp GetResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.URL != originalURL {
		t.Errorf("expected %q, got %q", originalURL, resp.URL)
	}
}