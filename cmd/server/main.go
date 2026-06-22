package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"typeahead-search-engine/internal/store"
	"typeahead-search-engine/internal/store/postgres"
	"typeahead-search-engine/internal/store/redis"
)

const (
	SERVER_PORT string = "localhost:8080"
)

// App represents the typeahead search application
type App struct {
	Store store.TypeaheadStore
}

func main() {
	storeType := os.Getenv("STORE_TYPE")
	var app App

	switch storeType {
	case "redis":
		app = App{Store: redis.New("localhost:6379")}
	case "postgres":
		app = App{Store: postgres.New("postgresql://postgres:admin123@localhost:5432/frequency-db")}
	default:
		fmt.Fprintf(os.Stderr, "server: invalid store type: %s. Use 'redis' or 'postgres'", storeType)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", app.healthCheckHandler)
	mux.HandleFunc("GET /api/v1/suggestions", app.suggestionsHandler)
	mux.HandleFunc("POST /api/v1/search", app.searchHandler)

	globalHandler := corsMiddleware(mux)

	fmt.Printf("Server is running on http://%s\n", SERVER_PORT)
	if err := http.ListenAndServe(SERVER_PORT, globalHandler); err != nil {
		fmt.Fprintf(os.Stderr, "error starting HTTP server: %v\n", err)
		os.Exit(1)
	}
}

func (app *App) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := w.Header()
		headers.Set("Access-Control-Allow-Origin", "*")
		headers.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		headers.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		next.ServeHTTP(w, r)
	})
}

func (app *App) suggestionsHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	prefix := params.Get("prefix")
	limit, _ := strconv.Atoi(params.Get("limit"))

	if len(prefix) == 0 || limit <= 0 {
		return
	}

	queries, err := app.Store.GetSuggestions(r.Context(), prefix, limit)
	if err != nil && err != context.Canceled {
		fmt.Fprintf(os.Stderr, "/suggestions: error fetching suggestions: %v\n", err)
		http.Error(w, "Error fetching suggestions", http.StatusInternalServerError)
		return
	}

	output, err := json.Marshal(queries)
	if err != nil {
		fmt.Fprintf(os.Stderr, "/suggestions: error converting queries into json: %v\n", err)
		http.Error(w, "Error transforming output", http.StatusInternalServerError)
		return
	}

	w.Write(output)
}

func (app *App) searchHandler(w http.ResponseWriter, r *http.Request) {
	// Read the search query from user.
	query, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "/search: error reading search query from body: %v\n", err)
		http.Error(w, "Error reading query", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Update frequncy for the provided query.
	app.Store.IncrementFrequency(r.Context(), string(query))

	// Notify client about resource creation.
	w.WriteHeader(http.StatusCreated)
}
