package storage

import (
	"context"
	"time"
)

type AStorage interface {
	Create(ctx context.Context, a []*Article) error
	GetByCategory(ctx context.Context, category string) ([]*Article, error)
	GetByAuthor(ctx context.Context, author string) ([]*Article, error)
	GetAll(ctx context.Context) ([]*Article, error)
}

type Article struct {
	ID        int       `db:"id"`
	ArticleID int       `db:"article_id"`
	Title     string    `db:"title"`
	Author    string    `db:"author"`
	Category  string    `db:"category"`
	URL       string    `db:"url"`
	CreatedAt time.Time `db:"created_at"`
}
