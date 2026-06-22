package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

// New creates and returns a new PostgreSQL store connection pool
func New(connString string) *PostgresStore {
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "postgres: error connecting to database: %v", err)
		os.Exit(1)
	}
	return &PostgresStore{pool: pool}
}

func (p *PostgresStore) GetSuggestions(ctx context.Context, prefix string, limit int) ([]string, error) {
	// Prepare prefix pattern for SQL LIKE query
	cleanPrefix := prefix + "%"
	sql := `
		SELECT query FROM search_queries
		WHERE query LIKE $1
		ORDER BY frequency DESC
		LIMIT $2;
	`

	// Query the database for matching suggestions
	rows, err := p.pool.Query(ctx, sql, cleanPrefix, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var query string
	queries := make([]string, 0, 10)

	// Iterate through results and collect suggestions
	for rows.Next() {
		rows.Scan(&query)
		queries = append(queries, query)
	}
	if err := rows.Err(); err != nil {
		return queries, err
	}

	return queries, nil
}

func (p *PostgresStore) IncrementFrequency(ctx context.Context, query string) error {
	sql := `
		SELECT frequency FROM search_queries
		WHERE query = $1
	`

	// Execute SQL query for fetching frequency of query.
	var freq int64
	row := p.pool.QueryRow(ctx, sql, query)
	row.Scan(&freq)

	// Define SQL query based on existence of query in database.
	if freq == 0 {
		sql = `
			INSERT INTO search_queries (query, frequency)
			VALUES ($1, $2);
		`
	} else {
		sql = `
			UPDATE search_queries
			SET frequency = $2
			WHERE query = $1;
		`
	}

	// Execute SQL query for incrementing the frequency of query by 1.
	if _, err := p.pool.Exec(ctx, sql, query, freq+1); err != nil {
		return err
	}
	return nil
}
