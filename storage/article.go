package storage

import (
	"context"
	"errors"
	"time"
)

type Storage interface {
	// article api
	Upsert(ctx context.Context, a []*ArticleAPI) error
	GetByCategoryAPI(ctx context.Context, category string) ([]*ArticleAPI, error)
	GetByAuthorAPI(ctx context.Context, userID int) ([]*ArticleAPI, error)
	GetAllAPI(ctx context.Context) ([]*ArticleAPI, error)
	GetCategoryAPI(ctx context.Context) ([]*ArticleCategoryDTO, error)
	GetAuthorAPI(ctx context.Context) ([]*ArticleAuthorDTO, error)
	GetLatest(ctx context.Context) ([]*ArticleAPI, error)
}

var ErrNoArticles = errors.New("no articles")

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

// ArticleAuthorDTO represents request artcile author object
type ArticleAuthorDTO struct {
	UserID        int    `json:"user_id" db:"user_id"`
	AuthorName    string `json:"first_name" db:"author_first_name"`
	AuthorSurname string `json:"last_name" db:"author_last_name"`
}

// ArticleCategoryDTO represents request DTO article category DTO
type ArticleCategoryDTO struct {
	Category string `db:"category"`
}
