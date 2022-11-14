package postgres

import (
	"context"

	"github.com/aibeksarsembayev/onelab-finalproject-telegrambot/storage"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	dbpool *sqlx.DB
}

func NewDBArticleRepo(dbpool *sqlx.DB) *Storage {
	return &Storage{
		dbpool: dbpool,
	}
}

// Upsert API articles ...
func (s *Storage) Upsert(ctx context.Context, a []*storage.ArticleAPI) error {
	_, err := s.dbpool.NamedExec(`INSERT INTO "article_api" (article_id, user_id, category_id, category, title, created_at, updated_at, author_first_name, author_last_name, url )
	VALUES (:article_id, :user_id, :category_id, :category, :title, :created_at, :updated_at, :author_first_name, :author_last_name, :url)
	ON CONFLICT(article_id) 
	DO UPDATE SET (user_id, category_id, category, title, created_at, updated_at, author_first_name, author_last_name, url) = 
	(excluded.user_id, excluded.category_id, excluded.category, excluded.title, excluded.created_at, excluded.updated_at, excluded.author_first_name, excluded.author_last_name, excluded.url)`, a)

	if err != nil {
		return err
	}
	return nil
}

// Create articles ...
func (s *Storage) Create(ctx context.Context, a []*storage.Article) error {
	_, err := s.dbpool.NamedExec(`INSERT INTO "article" (title, article_id, author, category, url, created_at)
	VALUES (:title, :article_id, :author, :category, :url, :created_at)
	ON CONFLICT(article_id) 
	DO UPDATE SET (title, author, category, url, created_at) = (excluded.title, excluded.author, excluded.category, excluded.url, excluded.created_at)`, a)

	if err != nil {
		return err
	}
	return nil
}

// GetByCategory ...
func (s *Storage) GetByCategory(ctx context.Context, category string) ([]*storage.Article, error) {
	a := []*storage.Article{}
	err := s.dbpool.Select(&a, `SELECT * FROM "article" WHERE category = $1 ORDER by created_at DESC`, category)
	if err != nil {
		return []*storage.Article{}, err
	}
	return a, nil
}

// GetByAuthor ...
func (s *Storage) GetByAuthor(ctx context.Context, author string) ([]*storage.Article, error) {
	a := []*storage.Article{}
	err := s.dbpool.Select(&a, `SELECT * FROM "article" WHERE author = $1 ORDER by created_at DESC`, author)
	if err != nil {
		return []*storage.Article{}, err
	}
	return a, nil
}

// GetAll ...
func (s *Storage) GetAll(ctx context.Context) ([]*storage.Article, error) {
	a := []*storage.Article{}
	err := s.dbpool.Select(&a, `SELECT * FROM "article" ORDER by created_at DESC`)
	if err != nil {
		return []*storage.Article{}, err
	}
	return a, nil
}

// GetCategory ...
func (s *Storage) GetCategory(ctx context.Context) ([]*storage.ArticleCategoryDTO, error) {
	c := []*storage.ArticleCategoryDTO{}
	err := s.dbpool.Select(&c, `SELECT category FROM "article" GROUP BY category`)
	if err != nil {
		return []*storage.ArticleCategoryDTO{}, err
	}
	return c, nil
}

// GetAuthor ...
func (s *Storage) GetAuthor(ctx context.Context) ([]*storage.ArticleAuthorDTO, error) {
	a := []*storage.ArticleAuthorDTO{}
	err := s.dbpool.Select(&a, `SELECT author FROM "article" GROUP BY author`)
	if err != nil {
		return []*storage.ArticleAuthorDTO{}, err
	}
	return a, nil
}
