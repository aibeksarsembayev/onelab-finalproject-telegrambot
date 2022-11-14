package apifetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/aibeksarsembayev/onelab-finalproject-telegrambot/storage"
	"github.com/aibeksarsembayev/onelab-finalproject-telegrambot/storage/postgres"
)

const (
	apiURL     = "https://sber-invest.kz/article/general/getall"
	articleURL = "https://sber-invest.kz/article/"
)

func NewFetcher(s *postgres.Storage, period time.Duration) {
	t := time.NewTicker(period * time.Minute)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			fetching(s)
		}
	}
}

func fetching(s *postgres.Storage) {

	url := apiURL

	client := http.Client{
		Timeout: 2 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// req.Header.Set("User-Agent","test")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	articles := []*storage.ArticleAPI{}
	err = json.Unmarshal(body, &articles)
	if err != nil {
		log.Fatal(err)
	}

	for _, article := range articles {
		article.URL = fmt.Sprintf("%v%v", articleURL, article.ArticleID)
	}

	err = s.Upsert(context.Background(), articles)
	if err != nil {
		log.Fatal("fetch api: can't insert or update fetched data into db: ", err)
	}

	log.Print("API fetching was done successfully")

	for _, v := range articles {
		fmt.Println(*v)
	}

}
