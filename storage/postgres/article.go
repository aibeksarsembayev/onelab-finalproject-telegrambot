package postgres

import (
	"context"
	"log"

	"github.com/aibeksarsembayev/onelab-finalproject-telegrambot/storage"
	parser "github.com/aibeksarsembayev/onelab-finalproject-telegrambot/tools/parser"
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

func (s *ArticleRepository) InsertArticle() {
	ax := parser.NewParser()

	err := s.Create(context.Background(), ax)
	if err != nil {
		log.Fatal("datavase insert issue", err)
	}
}

// Create articles ...
func (s *ArticleRepository) Create(ctx context.Context, a []*storage.Article) error {

	_, err := s.dbpool.NamedExec(`INSERT INTO "article" (title, article_id, author, category, url, created_at)
	VALUES (:title, :author, :category, :url, :created_at) 
	ON CONFLICT(article_id) 
	DO UPDATE SET title=excluded.title, author=excluded.author,  category=excluded.category, url=excluded.url, created_at=excluded.created_at`, a)

	if err != nil {
		return err
	}
	return nil
}

// GetByCategory ...
func (s *ArticleRepository) GetByCategory(ctx context.Context, category string) ([]*storage.Article, error) {
	return nil, nil
}

// GetByAuthor ...
func (s *ArticleRepository) GetByAuthor(ctx context.Context, author string) ([]*storage.Article, error) {
	return nil, nil
}

// GetAll ...
func (s *ArticleRepository) GetAll(ctx context.Context) ([]*storage.Article, error) {
	return nil, nil
}
