# Typeahead System

[![Typeahead Demo](./docs/demo.png)](./docs/demo.png)

A high-performance typeahead (search autocomplete) system designed to predict and suggest the most popular search queries as a user types. Built with a focus on system scalability, low-latency requirements, and efficient data ingestion for real-world search applications.

## Problem Statement

Design a typeahead application that:
- Provides real-time suggestions as users type.
- Shows the top 10 most relevant suggestions based on prefix matching and historical query frequency.
- Optimizes for **low latency** and **high availability**, even at the cost of eventual consistency.
- Handles a massive scale: ~2 million requests per second and ~320 TB of data over 20 years.

## Documentation

For a detailed breakdown of the requirements, design decisions, and data modeling, refer to the `docs/` directory:

- [01 Requirements and Scale](./docs/01-requirements-and-scale.md): Detailed functional/non-functional requirements and scale estimations.
- [02 High Level Design](./docs/02-high-level-design.md): Exploration of Trie vs. Key-Value store approaches, sharding strategies, and batching.
- [03 Dataset Information](./docs/03-dataset-information.md): Overview of the AOL search logs dataset used for this project.
- [04 ETL Pipeline Design](./docs/04-etl-pipeline-design.md): Architecture of the ingestion and loading pipeline.

## Tech Stack

- **Language:** Go (Golang)
- **Database:** PostgreSQL (using `pgx` driver)
- **Infrastructure:** Docker
- **Orchestration:** Makefile

## Local Setup

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- Make

### 1. Spin up the Database
Run a PostgreSQL instance in the background using Docker:

```bash
docker run -d \
  --name frequency-db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=admin123 \
  -e POSTGRES_DB=frequencydb \
  -p 5432:5432 \
  postgres:16
```

### 2. Configure Environment
Set the database connection URL as your environment variable:

```bash
export FREQUENCY_DB_URL="postgresql://postgres:admin123@localhost:5432/frequencydb"
```

### 3. Run the ETL Pipeline
The project uses a `Makefile` to manage the ingestion and loading of the AOL dataset.

```bash
# Run the full pipeline (Ingest -> Load)
make

# Or run individual steps:
make ingest    # Downloads and transforms the AOL dataset into data/dataset.csv
make load      # Performs a high-performance bulk copy of the CSV into Postgres
```

### 4. Inspect the Data
You can verify the loaded data using `psql`:

```bash
docker exec -it frequency-db psql -U postgres -d frequencydb -c "SELECT * FROM search_queries LIMIT 10;"
```

## Cleanup
To remove the generated dataset and stop the database:

```bash
make clean
docker stop frequency-db && docker rm frequency-db
```
