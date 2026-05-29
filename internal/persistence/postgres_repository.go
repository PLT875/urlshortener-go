package persistence

import (
	"database/sql"
	"fmt"

	"github.com/PLT875/urlshortener/internal/domain"
	_ "github.com/lib/pq"
)

var _ domain.Repository = (*PostgresRepository)(nil)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(dbURL string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	repo := &PostgresRepository{db: db}
	if err := repo.createTable(); err != nil {
		return nil, err
	}
	return repo, nil
}

func (p *PostgresRepository) createTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS urls (
		code VARCHAR(8) PRIMARY KEY,
		url TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := p.db.Exec(query)
	return err
}

func (p *PostgresRepository) Save(code, url string) {
	query := `INSERT INTO urls (code, url) VALUES ($1, $2) ON CONFLICT (code) DO NOTHING;`
	if err := p.db.QueryRow(query, code, url).Scan(); err != nil && err != sql.ErrNoRows {
		fmt.Printf("error saving URL: %v\n", err)
	}
}

func (p *PostgresRepository) Get(code string) (string, bool) {
	query := `SELECT url FROM urls WHERE code = $1;`
	var url string
	err := p.db.QueryRow(query, code).Scan(&url)
	if err == sql.ErrNoRows {
		return "", false
	}
	if err != nil {
		fmt.Printf("error retrieving URL: %v\n", err)
		return "", false
	}
	return url, true
}

func (p *PostgresRepository) Close() error {
	return p.db.Close()
}
