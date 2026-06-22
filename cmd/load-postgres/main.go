package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	// Establish connection to PostgreSQL database
	conn, err := pgx.Connect(context.Background(), `postgresql://postgres:admin123@localhost:5432/frequency-db`)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	// Initialize search queries table
	tableName := "search_queries"
	columnNames := []string{"query", "frequency"}
	createTableQuery := fmt.Sprintf(
		`
		CREATE TABLE IF NOT EXISTS %s (
			%s TEXT PRIMARY KEY,
			%s BIGINT NOT NULL DEFAULT 0
		);
		`,
		tableName, columnNames[0], columnNames[1],
	)
	if _, err := conn.Exec(context.Background(), createTableQuery); err != nil {
		log.Fatalf("Unable to create table: %s: %v", tableName, err)
	}

	// Open dataset CSV file
	f, err := os.Open("data/dataset.csv")
	if err != nil {
		log.Fatalf("Unable to open dataset: %v", err)
	}

	// Prepare bulk copy command, HEADER true tells PostgreSQL
	// to ignore the first header line consisting of column names.
	bulkCopyQuery := fmt.Sprintf(
		"COPY %s (%s, %s) FROM STDIN WITH (FORMAT csv, HEADER true)",
		tableName, columnNames[0], columnNames[1],
	)

	// Bulk load dataset into PostgreSQL
	if _, err := conn.PgConn().CopyFrom(context.Background(), f, bulkCopyQuery); err != nil {
		log.Fatalf("Bulk copy failed: %v", err)
	}
}
