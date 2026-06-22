package store

import "context"

type TypeaheadStore interface {
	GetSuggestions(ctx context.Context, prefix string, limit int) ([]string, error)
	IncrementFrequency(ctx context.Context, query string) error
}
