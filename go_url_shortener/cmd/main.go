package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"

	"github.com/vktnwar/go_url_shortener/config"
	"github.com/vktnwar/go_url_shortener/repository"
	"github.com/vktnwar/go_url_shortener/server"
)

func main() {
	// Carrega configurações
	cfg := config.LoadConfig()

	// =======================
	// Conexão Redis
	// =======================
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
		DB:       cfg.RedisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Erro ao conectar ao Redis: %v", err)
	}
	log.Println("✅ Redis conectado")
	redisRepo := repository.NewRedisRepository(rdb)

	// =======================
	// Conexão Postgres
	// =======================
	db, err := sql.Open("postgres", cfg.PostgresURL)
	if err != nil {
		log.Fatalf("Erro ao conectar ao Postgres: %v", err)
	}

	// Testa a conexão
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Erro ao testar conexão com Postgres: %v", err)
	}

	log.Println("✅ Postgres conectado")
	pgRepo := repository.NewPostgresRepository(db)

	// =======================
	// Start do servidor
	// =======================
	r := server.NewRouter(pgRepo, redisRepo, cfg)

	log.Println("🚀 Servidor rodando em :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
