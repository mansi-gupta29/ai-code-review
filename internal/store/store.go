package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mansi-gupta29/ai-code-review/internal/models"
)

type Store struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

// SaveReview inserts a review and returns it with the generated ID + timestamp
func (s *Store) SaveReview(ctx context.Context, r models.Review) (models.Review, error) {
	query := `
		INSERT INTO reviews (language, code, review)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	err := s.pool.QueryRow(ctx, query, r.Language, r.Code, r.Review).
		Scan(&r.ID, &r.CreatedAt)
	if err != nil {
		return models.Review{}, err
	}
	return r, nil
}

// ListReviews returns all saved reviews, newest first
func (s *Store) ListReviews(ctx context.Context) ([]models.Review, error) {
	query := `
		SELECT id, language, code, review, created_at
		FROM reviews
		ORDER BY created_at DESC`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var r models.Review
		if err := rows.Scan(&r.ID, &r.Language, &r.Code, &r.Review, &r.CreatedAt); err != nil {
			return nil, err
		}
		reviews = append(reviews, r)
	}
	return reviews, rows.Err()
}

func (s *Store) GetReviewByID(ctx context.Context, id int) (models.Review, error) {
	query := `
        SELECT id, language, code, review, created_at
        FROM reviews
        WHERE id = $1`

	var r models.Review
	err := s.pool.QueryRow(ctx, query, id).
		Scan(&r.ID, &r.Language, &r.Code, &r.Review, &r.CreatedAt)
	if err != nil {
		return models.Review{}, err
	}

	return r, nil
}
