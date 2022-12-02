package apifetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/zecodein/sber-invest-bot/storage"
	"github.com/zecodein/sber-invest-bot/storage/postgres"
	"go.uber.org/zap"
)

const (
	apiURL     = "https://sber-invest.kz/article/general/getall"
	articleURL = "https://sber-invest.kz/article/"
)

type Fetcher struct {
	lg      *zap.Logger
	storage *postgres.Storage
	period  time.Duration
}

func New(logger *zap.Logger, storage *postgres.Storage, period time.Duration) *Fetcher {
	// initial fetch
	return &Fetcher{
		lg:      logger,
		storage: storage,
		period:  period,
	}
}

func (f *Fetcher) Fetch() {
	// initial fetch
	err := f.fetching()
	if err != nil {
		f.lg.Sugar().Error(err)
	}
	t := time.NewTicker(f.period * time.Minute)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			err = f.fetching()
			if err != nil {
				f.lg.Sugar().Error(err)
			}
		}
	}
}

func (f *Fetcher) fetching() error {
	url := apiURL

	client := http.Client{
		Timeout: 2 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("fetch api: can't wrap request: %w", err)
	}

	// req.Header.Set("User-Agent","test")

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("fetch api: can't send request: %w", err)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("fetch api: can't read response: %w", err)
	}

	articles := []*storage.ArticleAPI{}
	err = json.Unmarshal(body, &articles)
	if err != nil {
		return fmt.Errorf("fetch api: can't unmarshal response: %w", err)
	}

	for _, article := range articles {
		article.URL = fmt.Sprintf("%v%v", articleURL, article.ArticleID)
	}

	err = f.storage.Upsert(context.Background(), articles)
	if err != nil {
		// log.Fatal("fetch api: can't insert or update fetched data into db: ", err)
		return fmt.Errorf("fetch api: can't insert or update fetched data into db: %w", err)
	}

	f.lg.Info("API fetching was done successfully")
	return nil
}
