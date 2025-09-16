package repository

import (
	"context"
	"database/sql"
)

type PostgresRepository struct {
	DB *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{DB: db}
}

// SaveURL insere uma nova URL encurtada
func (r *PostgresRepository) SaveURL(ctx context.Context, shortID, originalURL string) error {
	query := `
		INSERT INTO urls (short_id, original_url, clicks)
		VALUES ($1, $2, 0)
		ON CONFLICT (short_id) DO NOTHING;
	`
	_, err := r.DB.ExecContext(ctx, query, shortID, originalURL)
	return err
}

// GetOriginalURL busca a URL original pelo shortID
func (r *PostgresRepository) GetOriginalURL(ctx context.Context, shortID string) (string, error) {
	var originalURL string
	query := `SELECT original_url FROM urls WHERE short_id = $1`
	err := r.DB.QueryRowContext(ctx, query, shortID).Scan(&originalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return originalURL, nil
}

// IncrementClicks incrementa o contador de cliques
func (r *PostgresRepository) IncrementClicks(ctx context.Context, shortID string) error {
	query := `UPDATE urls SET clicks = clicks + 1 WHERE short_id = $1`
	_, err := r.DB.ExecContext(ctx, query, shortID)
	return err
}
