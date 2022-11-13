package parser

import (
	"context"
	"fmt"
	"log"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/aibeksarsembayev/onelab-finalproject-telegrambot/storage"
	"github.com/aibeksarsembayev/onelab-finalproject-telegrambot/storage/postgres"
	"github.com/gocolly/colly"
)

func NewParser(s *postgres.AStorage, period time.Duration) {
	t := time.NewTicker(period * time.Minute)
	defer t.Stop()
	for {
		select {
		case <-t.C: // initiate parsing
			parsing(s)
		}
	}

}

func parsing(s *postgres.AStorage) {
	ax := []*storage.Article{}

	for i := 1; i < 8; i++ {
		tempax := parseByCat(i)
		ax = append(ax, tempax...)
	}

	err := s.Create(context.TODO(), ax)
	if err != nil {
		log.Fatal("parse: can't insert parsed data into db", err)
	}
	log.Print("periodical parsing was done successfully")
}

func parseByCat(catN int) []*storage.Article {
	c := colly.NewCollector(
		colly.AllowedDomains("sber-invest.kz"),
	)
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	ax := []*storage.Article{}

	cat := ""

	c.OnHTML("header.header", func(e *colly.HTMLElement) {
		cat = e.ChildText("h1.heading")
	})

	c.OnHTML("div.message", func(e *colly.HTMLElement) {
		a := &storage.Article{}

		a.Title = e.ChildText("#article-title")
		a.Author = e.ChildText("#article-author-name")

		timeStr := e.ChildText("#article-create-date")

		a.CreatedAt = parseTime(timeStr)

		a.Category = cat
		a.URL = "https://sber-invest.kz" + e.ChildAttr("a#article-title", "href")

		articleIDstr := path.Base(a.URL)
		a.ArticleID, _ = strconv.Atoi(articleIDstr)

		if a.Author != "" {
			ax = append(ax, a)
		}
	})

	url := fmt.Sprintf("https://sber-invest.kz/knowledgebase/%d", catN)

	c.Visit(url)

	return ax
}

func parseTime(timeStr string) time.Time {

	if timeStr == "" {
		return time.Time{}
	}

	timeStr = strings.TrimPrefix(timeStr, "Добавлен: ")

	yy := timeStr[7:9]
	mm := timeStr[3:6]
	dd := timeStr[:2]
	tt := timeStr[10:]

	timeFormatted := fmt.Sprintf("20%s-%s-%sT%s", yy, mm, dd, tt)
	timeCreated, err := time.Parse("2006-Jan-02T15:04", timeFormatted)

	if err != nil {
		log.Fatal("time parse issue", err)
	}

	return timeCreated
}
