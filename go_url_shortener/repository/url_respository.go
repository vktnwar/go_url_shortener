package repository

import (
	"database/sql"

	"github.com/vktnwar/go_url_shortener/models"
)

type PostgresURLRepository struct {
	DB *sql.DB
}

func NewPostgresURLRepository(db *sql.DB) *PostgresURLRepository {
	return &PostgresURLRepository{DB: db}
}

func (r *PostgresURLRepository) Save(url *models.URL) error {
	query := `INSERT INTO urls (original, short, clicks, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	return r.DB.QueryRow(query, url.Original, url.Short, url.Clicks, url.CreatedAt).Scan(&url.ID)
}

func (r *PostgresURLRepository) FindByShort(short string) (*models.URL, error) {
	query := `SELECT id, original, short, clicks, created_at FROM urls WHERE short=$1`
	row := r.DB.QueryRow(query, short)

	var u models.URL
	err := row.Scan(&u.ID, &u.Original, &u.Short, &u.Clicks, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PostgresURLRepository) IncrementClicks(short string) error {
	_, err := r.DB.Exec(`UPDATE urls SET clicks = clicks + 1 WHERE short=$1`, short)
	return err
}
