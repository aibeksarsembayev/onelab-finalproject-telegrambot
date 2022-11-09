package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type ArticleRepository struct {
	dbpool *sqlx.DB
}

func NewDBArticleRepo(dbpool *sqlx.DB) *ArticleRepository {
	return &ArticleRepository{
		dbpool: dbpool,
	}
}

func (s *ArticleRepository) Create(ctx context.Context) error {
	return nil
}
