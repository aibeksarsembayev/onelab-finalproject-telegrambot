package storage

import (
	"context"
	"errors"
	"time"
)

type Storage interface {
	Create(ctx context.Context, a []*Article) error
	GetByCategory(ctx context.Context, category string) ([]*Article, error)
	GetByAuthor(ctx context.Context, author string) ([]*Article, error)
	GetAll(ctx context.Context) ([]*Article, error)
	GetCategory(ctx context.Context) ([]*ArticleCategoryDTO, error)
	GetAuthor(ctx context.Context) ([]*ArticleAuthorDTO, error)
}

var ErrNoArticles = errors.New("no articles")

// Article represents article object
type Article struct {
	ID        int       `db:"id"`
	ArticleID int       `db:"article_id"`
	Title     string    `db:"title"`
	Author    string    `db:"author"`
	Category  string    `db:"category"`
	URL       string    `db:"url"`
	CreatedAt time.Time `db:"created_at"`
}

// ArticleAuthorDTO represents request artcile author object
type ArticleAuthorDTO struct {
	Author string `db:"author"`
}

// ArticleCategoryDTO represents request DTO article category DTO
type ArticleCategoryDTO struct {
	Category string `db:"category"`
}

// ArticleAPI represents article to fetch from API
type ArticleAPI struct {
	ID            int       `json:"article_api_id" db:"id"`
	ArticleID     int       `json:"id" db:"article_id"`
	UserID        int       `json:"user_id" db:"user_id"`
	CategoryID    int       `json:"category_id" db:"category_id"`
	Category      string    `json:"category_name" db:"category"`
	Title         string    `json:"title" db:"title"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	AuthorName    string    `json:"first_name" db:"author_first_name"`
	AuthorSurname string    `json:"last_name" db:"author_last_name"`
	URL           string    `json:"url" db:"url"`
}
