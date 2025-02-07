// internal/adapters/secondary/file/repository.go
package file

import (
	"act/pkg/track/domain/rating"
	"act/pkg/track/ports/secondary"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileRepository struct {
	mu       sync.RWMutex
	filepath string
	ratings  map[string]rating.DayRating
}

func NewFileRepository(path string) (secondary.RatingRepository, error) {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("creating directory: %w", err)
	}

	repo := &FileRepository{
		filepath: path,
		ratings:  make(map[string]rating.DayRating),
	}

	// Load existing data if file exists
	if err := repo.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("loading ratings: %w", err)
	}

	return repo, nil
}

func (r *FileRepository) load() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := os.ReadFile(r.filepath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &r.ratings)
}

func (r *FileRepository) save() error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data, err := json.MarshalIndent(r.ratings, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling ratings: %w", err)
	}

	return os.WriteFile(r.filepath, data, 0644)
}

func (r *FileRepository) Save(_ context.Context, dr rating.DayRating) error {
	r.mu.Lock()
	r.ratings[dr.ID] = dr
	r.mu.Unlock()

	return r.save()
}

func (r *FileRepository) GetByID(_ context.Context, id string) (rating.DayRating, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	dr, exists := r.ratings[id]
	if !exists {
		return rating.DayRating{}, rating.ErrNotFound
	}

	return dr, nil
}

func (r *FileRepository) GetByDateRange(_ context.Context, start, end time.Time) ([]rating.DayRating, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []rating.DayRating
	for _, dr := range r.ratings {
		if (dr.Date.Equal(start) || dr.Date.After(start)) &&
			(dr.Date.Equal(end) || dr.Date.Before(end)) {
			results = append(results, dr)
		}
	}

	return results, nil
}

func (r *FileRepository) GetByWeek(_ context.Context, year, week int) ([]rating.DayRating, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []rating.DayRating
	for _, dr := range r.ratings {
		y, w := dr.Date.ISOWeek()
		if y == year && w == week {
			results = append(results, dr)
		}
	}

	return results, nil
}
