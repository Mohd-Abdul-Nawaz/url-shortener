package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"url-shortener/internal/store"
)

type ShortenRequest struct {
	URL         string `json:"url"`
	CustomAlias string `json:"custom_alias,omitempty"`
	ExpiresIn   int    `json:"expires_in,omitempty"` // seconds, 0 = never
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
	Alias    string `json:"alias"`
	ExpiresAt string `json:"expires_at,omitempty"`
}

type StatsResponse struct {
	Alias      string `json:"alias"`
	OriginalURL string `json:"original_url"`
	ClickCount  int    `json:"click_count"`
	CreatedAt   string `json:"created_at"`
}

type ShortenerHandler struct {
	store store.Store
}

func NewShortenerHandler(s store.Store) *ShortenerHandler {
	return &ShortenerHandler{store: s}
}

func (h *ShortenerHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate URL
	_, err := url.ParseRequestURI(req.URL)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	alias := req.CustomAlias
	if alias == "" {
		alias = generateAlias(7)
	} else if !isValidAlias(alias) {
		http.Error(w, "Invalid alias format (alphanumeric, dash, underscore only)", http.StatusBadRequest)
		return
	}

	// Check if alias already exists
	if h.store.Exists(alias) {
		http.Error(w, "Alias already taken", http.StatusConflict)
		return
	}

	err = h.store.Set(alias, req.URL, req.ExpiresIn)
	if err != nil {
		http.Error(w, "Failed to save URL", http.StatusInternalServerError)
		return
	}

	baseURL := getBaseURL(r)
	resp := ShortenResponse{
		ShortURL:  baseURL + "/" + alias,
		Alias:     alias,
		ExpiresAt: "", // populated by store if needed
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *ShortenerHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	alias := strings.TrimPrefix(r.URL.Path, "/")
	if alias == "" || alias == "favicon.ico" {
		http.Redirect(w, r, "/docs", http.StatusFound)
		return
	}

	originalURL, err := h.store.Get(alias)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	h.store.IncrementClick(alias)
	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}

func generateAlias(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[hashByte(i)]
	}
	return string(result)
}

func hashByte(i int) byte {
	// Simple pseudo-random for demo; use crypto/rand in production
	return byte((i*31 + 17) % 62)
}

func isValidAlias(alias string) bool {
	for _, c := range alias {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') || c == '-' || c == '_') {
			return false
		}
	}
	return len(alias) > 0 && len(alias) <= 50
}

func getBaseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + r.Host
}