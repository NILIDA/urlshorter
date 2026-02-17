package postgres

import (
	"database/sql"
	"errors"

	"urlshort/internal/generator"
	"urlshort/internal/storage"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	ErrGenerateShort = errors.New("failed to generate short url")
)

type PostgresStorage struct{
	db *sql.DB
}

func NewPostgresStorage(connStr string) (*PostgresStorage, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil{
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	if err := createTable(db); err != nil{
		return nil, err
	}
	return &PostgresStorage{db: db}, nil
}


func createTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS urls(
	id SERIAL PRIMARY KEY,
	short_url VARCHAR(10) UNIQUE NOT NULL,
	origin_url TEXT	UNIQUE NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_short_url ON urls (short_url);
    CREATE INDEX IF NOT EXISTS idx_origin_url ON urls (origin_url);
	`
	_, err := db.Exec(query)
	return err
}

func (p *PostgresStorage) Save(origURL string) (string, error) {
	var existingShort string
	err := p.db.QueryRow("SELECT short_url FROM urls WHERE origin_url = $1", origURL).Scan(&existingShort)
	if err == nil{
		return existingShort, nil
	} else if !errors.Is(err, sql.ErrNoRows){
		return "", err
	}

	for i := 0; i < 10; i++ {
		short := generator.Generate()
		res, err := p.db.Exec(
			"INSERT INTO urls (short_url, origin_url) VALUES ($1, $2) ON CONFLICT (short_url) DO NOTHING",
			short, origURL,
		)
		if err != nil {
			return "", err
		}
		rowsAffected, _ := res.RowsAffected()
		if rowsAffected > 0 {
			return short, nil
		}
	}
	return "", ErrGenerateShort
}

func (p *PostgresStorage) Get(shortURL string) (string, error) {
	var original string
	err := p.db.QueryRow("SELECT origin_url FROM urls WHERE short_url = $1", shortURL).Scan(&original)
	if err != nil{
		if errors.Is(err, sql.ErrNoRows){
			return "", storage.ErrNotFound
		}
		return "", err
	}
	return original, nil
}

func (p *PostgresStorage) Close() error {
	return p.db.Close()
}