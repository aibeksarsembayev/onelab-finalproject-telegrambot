package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zecodein/sber-invest-bot/storage"
	randomizer "github.com/zecodein/sber-invest-bot/tools/random"
)

func TestStorage_Upsert(t *testing.T) {
	args := []*storage.ArticleAPI{}
	for i := 0; i < 10; i++ {
		arg := randomArticle()
		args = append(args, arg)
	}
	err := testStorage.Upsert(context.Background(), args)
	require.NoError(t, err)
}

func TestStorage_GetByCategoryAPI(t *testing.T) {
	articles := []*storage.ArticleAPI{}
	category, err := testStorage.GetCategoryAPI(context.Background())
	require.NoError(t, err)
	for i := 0; i < len(category); i++ {
		articles, err = testStorage.GetByCategoryAPI(context.Background(), category[i].Category)
	}

	require.NoError(t, err)

	for _, article := range articles {
		require.NotEmpty(t, article)
	}
}

func TestStorage_GetByAuthorAPI(t *testing.T) {
	articles := []*storage.ArticleAPI{}
	authors, err := testStorage.GetAuthorAPI(context.Background())
	require.NoError(t, err)
	for i := 0; i < len(authors); i++ {
		author := authors[i].AuthorName + " " + authors[i].AuthorSurname
		articles, err = testStorage.GetByCategoryAPI(context.Background(), author)
	}

	require.NoError(t, err)

	for _, article := range articles {
		require.NotEmpty(t, article)
	}
}

func TestStorage_GetAllAPI(t *testing.T) {
	articles, err := testStorage.GetAllAPI(context.Background())

	require.NoError(t, err)

	for _, article := range articles {
		require.NotEmpty(t, article)
	}
}

func TestStorage_GetCategoryAPI(t *testing.T) {
	category, err := testStorage.GetCategoryAPI(context.Background())
	require.NoError(t, err)

	for _, c := range category {
		require.NotEmpty(t, c)
	}
}

func TestStorage_GetLatest(t *testing.T) {
	articles, err := testStorage.GetLatest(context.Background())

	require.NoError(t, err)

	for _, article := range articles {
		require.NotEmpty(t, article)
	}
}

func TestStorage_GetAuthorAPI(t *testing.T) {
	author, err := testStorage.GetCategoryAPI(context.Background())
	require.NoError(t, err)

	for _, a := range author {
		require.NotEmpty(t, a)
	}
}

func randomArticle() *storage.ArticleAPI {
	return &storage.ArticleAPI{
		ArticleID:     int(randomizer.RandomInt(1, 999)),
		UserID:        int(randomizer.RandomInt(1, 9)),
		CategoryID:    int(randomizer.RandomInt(1, 9)),
		Category:      randomizer.RandomString(10),
		Title:         randomizer.RandomString(10),
		CreatedAt:     randomizer.RandomDate(),
		UpdatedAt:     randomizer.RandomDate().Add(10),
		AuthorName:    randomizer.RandomString(6),
		AuthorSurname: randomizer.RandomString(10),
		URL:           randomizer.RandomString(20),
	}
}
