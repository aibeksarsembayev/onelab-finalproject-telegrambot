package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/zecodein/sber-invest-bot/storage"
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

// GetByCategoryAPI ...
func (s *Storage) GetByCategoryAPI(ctx context.Context, category string) ([]*storage.ArticleAPI, error) {
	a := []*storage.ArticleAPI{}
	err := s.dbpool.Select(&a, `SELECT * FROM "article_api" WHERE category = $1 ORDER by created_at DESC`, category)
	if err != nil {
		return []*storage.ArticleAPI{}, err
	}
	return a, nil
}

// GetByAuthorAPI ...
func (s *Storage) GetByAuthorAPI(ctx context.Context, userID int) ([]*storage.ArticleAPI, error) {
	a := []*storage.ArticleAPI{}
	err := s.dbpool.Select(&a, `SELECT * FROM "article_api" WHERE user_id = $1 ORDER by created_at DESC`, userID)
	if err != nil {
		return []*storage.ArticleAPI{}, err
	}
	return a, nil
}

// GetAllAPI ...
func (s *Storage) GetAllAPI(ctx context.Context) ([]*storage.ArticleAPI, error) {
	a := []*storage.ArticleAPI{}
	err := s.dbpool.Select(&a, `SELECT * FROM "article_api" ORDER by created_at DESC`)
	if err != nil {
		return []*storage.ArticleAPI{}, err
	}
	return a, nil
}

// GetCategoryAPI ...
func (s *Storage) GetCategoryAPI(ctx context.Context) ([]*storage.ArticleCategoryDTO, error) {
	c := []*storage.ArticleCategoryDTO{}
	err := s.dbpool.Select(&c, `SELECT category FROM "article_api" GROUP BY category`)
	if err != nil {
		return []*storage.ArticleCategoryDTO{}, err
	}
	return c, nil
}

// GetAuthorAPI ...
func (s *Storage) GetAuthorAPI(ctx context.Context) ([]*storage.ArticleAuthorDTO, error) {
	a := []*storage.ArticleAuthorDTO{}
	err := s.dbpool.Select(&a, `SELECT user_id, author_first_name, author_last_name FROM "article_api" GROUP BY user_id, author_first_name, author_last_name`)
	if err != nil {
		return []*storage.ArticleAuthorDTO{}, err
	}
	return a, nil
}

// GetLatest article for 7 days ...
func (s *Storage) GetLatest(ctx context.Context) ([]*storage.ArticleAPI, error) {
	a := []*storage.ArticleAPI{}
	err := s.dbpool.Select(&a, `SELECT * FROM "article_api" 
	WHERE created_at > current_date - interval '7' day
	ORDER by created_at DESC`)
	if err != nil {
		return []*storage.ArticleAPI{}, err
	}
	return a, nil
}
