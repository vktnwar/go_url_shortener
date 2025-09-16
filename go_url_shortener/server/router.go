package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/vktnwar/go_url_shortener/config"
	"github.com/vktnwar/go_url_shortener/middleware"
	"github.com/vktnwar/go_url_shortener/repository"
	"github.com/vktnwar/go_url_shortener/service"
)

type URLHandler struct {
	Service *service.URLService
}

func NewRouter(pgRepo *repository.PostgresRepository, redisRepo *repository.RedisRepository, cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	urlService := service.NewURLService(pgRepo, redisRepo)
	handler := &URLHandler{Service: urlService}

	// Healthcheck
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Endpoint protegido com Rate Limiting (5 req/min por IP)
	r.With(middleware.RateLimiterMiddleware(redisRepo.Client, 5, time.Minute)).
		Post("/shorten", handler.Shorten)

	// Redirect sem limite
	r.Get("/{shortID}", handler.Resolve)

	return r
}

// ===== Handlers =====

// Shorten gera uma nova URL encurtada
func (h *URLHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	type request struct {
		URL string `json:"url"`
	}
	type response struct {
		ShortURL string `json:"short_url"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	shortID, err := h.Service.ShortenURL(ctx, req.URL)
	if err != nil {
		http.Error(w, "could not shorten URL", http.StatusInternalServerError)
		return
	}

	resp := response{ShortURL: "http://localhost:8080/" + shortID}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Resolve redireciona uma shortURL para a original
func (h *URLHandler) Resolve(w http.ResponseWriter, r *http.Request) {
	shortID := chi.URLParam(r, "shortID")

	ctx := context.Background()
	originalURL, err := h.Service.ResolveURL(ctx, shortID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}
