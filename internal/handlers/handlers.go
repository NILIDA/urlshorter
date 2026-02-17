package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"urlshort/internal/storage"
)

type SaveRequest struct {
	URL string `json:"url"`
}

type SaveResponse struct {
	Short string `json:"short"`
}

type GetResponse struct {
	URL string `json:"url"`
}

type Handler struct {
	store   storage.Storage
	logger  *zap.SugaredLogger
	baseURL string
	tmpl    *template.Template
}

func NewHandler(store storage.Storage, logger *zap.SugaredLogger, baseURL string) *Handler {
	tmpl, err := template.ParseGlob("internal/handlers/web/templates/*.html")
	if err != nil {
		logger.Fatalw("failed to load template", "error", err)
	}

	return &Handler{
		store:   store,
		logger:  logger,
		baseURL: baseURL,
		tmpl:    tmpl,
	}
}

type PageData struct {
	InputURL string
	Result   string
	Error    string
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	data := PageData{}
	h.tmpl.Execute(w, data)
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	inputURL := r.FormValue("url")
	action := r.FormValue("action")

	data := PageData{InputURL: inputURL}

	if inputURL == "" {
		data.Error = "Пожалуйста, введите ссылку или код"
		h.tmpl.Execute(w, data)
		return
	}

	if action == "shorten" {
		if !strings.HasPrefix(inputURL, "http://") && !strings.HasPrefix(inputURL, "https://") {
			data.Error = "Ссылка должна начинаться с http:// или https://"
			h.tmpl.Execute(w, data)
			return
		}

		short, err := h.store.Save(inputURL)
		if err != nil {
			h.logger.Errorw("failed to save url", "error", err)
			data.Error = "Ошибка при сохранении ссылки"
			h.tmpl.Execute(w, data)
			return
		}

		data.Result = short

	} else if action == "expand" {
		shortCode := inputURL
		if len(inputURL) > len(h.baseURL) && inputURL[:len(h.baseURL)] == h.baseURL {
			shortCode = inputURL[len(h.baseURL)+1:]
		}

		if len(shortCode) != 10 {
			data.Error = "Короткий код должен быть длиной 10 символов"
			h.tmpl.Execute(w, data)
			return
		}

		longURL, err := h.store.Get(shortCode)
		if err != nil {
			if err == storage.ErrNotFound {
				data.Error = "Ссылка не найдена"
			} else {
				h.logger.Errorw("failed to get url", "error", err)
				data.Error = "Ошибка при получении ссылки"
			}
			h.tmpl.Execute(w, data)
			return
		}

		data.Result = longURL
	}

	h.tmpl.Execute(w, data)
}

func (h *Handler) Save(w http.ResponseWriter, r *http.Request) {
	var req SaveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorw("failed to decode request", "error", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	_, err := url.ParseRequestURI(req.URL)
	if err != nil {
		h.logger.Infow("invalid url provided", "url", req.URL)
		http.Error(w, "invalid url format", http.StatusBadRequest)
		return
	}

	short, err := h.store.Save(req.URL)
	if err != nil {
		h.logger.Errorw("failed to save url", "url", req.URL, "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	resp := SaveResponse{Short: short}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Errorw("failed to encode response", "error", err)
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	short := vars["short"]

	original, err := h.store.Get(short)
	if err != nil {
		if err == storage.ErrNotFound {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		h.logger.Errorw("failed to get url", "short", short, "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("Accept") == "application/json" {
		resp := GetResponse{URL: original}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			h.logger.Errorw("failed to encode response", "error", err)
		}
		return
	}

	http.Redirect(w, r, original, http.StatusFound)
}
