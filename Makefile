.PHONY: all ingest load-postgres load-redis clean

# Typeahead Search Engine - Build and data loading pipeline
all: load

data/dataset.csv: cmd/ingest/main.go
	@echo "==> Extracting and Transforming data..."
	@mkdir -p data
	go run cmd/ingest/main.go

load-postgres: data/dataset.csv
	@echo "==> Loading data into PostgreSQL..."
	go run cmd/load-postgres/main.go

load-redis: data/dataset.csv
	@echo "==> Loading data into Redis..."
	go run cmd/load-redis/main.go

clean:
	@echo "==> Cleaning up artifacts..."
	rm -f data/dataset.csv
