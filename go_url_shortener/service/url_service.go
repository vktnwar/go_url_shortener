package service

import (
	"context"
	"errors"

	"github.com/segmentio/ksuid"
	"github.com/vktnwar/go_url_shortener/repository"
)

type URLService struct {
	PgRepo    *repository.PostgresRepository
	RedisRepo *repository.RedisRepository
}

func NewURLService(pgRepo *repository.PostgresRepository, redisRepo *repository.RedisRepository) *URLService {
	return &URLService{
		PgRepo:    pgRepo,
		RedisRepo: redisRepo,
	}
}

// ShortenURL gera um shortID e salva no Postgres
func (s *URLService) ShortenURL(ctx context.Context, originalURL string) (string, error) {
	shortID := ksuid.New().String()[:8] // gera 8 caracteres aleatórios

	// Salva no Postgres
	if err := s.PgRepo.SaveURL(ctx, shortID, originalURL); err != nil {
		return "", err
	}

	// Salva no Redis cache (opcional)
	_ = s.RedisRepo.Client.Set(ctx, shortID, originalURL, 0).Err()

	return shortID, nil
}

// ResolveURL busca a URL original pelo shortID
func (s *URLService) ResolveURL(ctx context.Context, shortID string) (string, error) {
	// 1️⃣ Tenta pegar do Redis
	originalURL, err := s.RedisRepo.Client.Get(ctx, shortID).Result()
	if err == nil {
		// Se encontrou no cache, incrementa clique no Postgres
		_ = s.PgRepo.IncrementClicks(ctx, shortID)
		return originalURL, nil
	}

	// 2️⃣ Se não encontrou no cache, busca no Postgres
	originalURL, err = s.PgRepo.GetOriginalURL(ctx, shortID)
	if err != nil {
		return "", err
	}
	if originalURL == "" {
		return "", errors.New("URL não encontrada")
	}

	// Atualiza o cache do Redis e incrementa clique
	_ = s.RedisRepo.Client.Set(ctx, shortID, originalURL, 0).Err()
	_ = s.PgRepo.IncrementClicks(ctx, shortID)

	return originalURL, nil
}
