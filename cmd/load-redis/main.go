package main

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		Protocol: 2,
		DB:       0,
	})

	fmt.Println("Redis store initialized and ready for data ingestion")

	f, err := os.Open("data/dataset.csv")
	if err != nil {
		log.Fatalf("Unable to open dataset: %v", err)
	}

	r := csv.NewReader(f)
	ctx := context.Background()
	cnt := 0
	for {
		record, err := r.Read()

		if errors.Is(err, io.EOF) {
			fmt.Println("All records have been processed.")
			return
		}
		if err != nil {
			fmt.Printf("Error reading record: %v\n", err)
			continue
		}

		query := record[0]
		freq, _ := strconv.Atoi(record[1])

		// Index all prefixes of the query in Redis
		for i := range query {
			if err := rdb.ZAdd(ctx, query[:i+1], redis.Z{Score: float64(freq), Member: query}).Err(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to insert: %q: %v\n", query[:i+1], err)
			}
		}

		cnt++
		if cnt%1000 == 0 {
			fmt.Printf("Queries processed: %d%%		\r", cnt)
			os.Stdout.Sync()
		}
	}
}
