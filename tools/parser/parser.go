package parser

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type Article struct {
	Title string
	// Content   string
	Author    string
	Category  string
	URL       string
	CreatedAt time.Time
}

func NewParser() {
	ax := []*Article{}

	for i := 1; i < 8; i++ {
		tempax := parseByCat(i)
		ax = append(ax, tempax...)
	}

	for i, a := range ax {
		fmt.Printf("========== #%d ==========\n%s\n", i, a)
	}

}

func parseByCat(catN int) []*Article {
	c := colly.NewCollector(
		colly.AllowedDomains("sber-invest.kz"),
	)
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	ax := []*Article{}

	cat := ""

	c.OnHTML("header.header", func(e *colly.HTMLElement) {
		cat = e.ChildText("h1.heading")
	})

	c.OnHTML("div.message", func(e *colly.HTMLElement) {
		a := &Article{}

		a.Title = e.ChildText("#article-title")
		a.Author = e.ChildText("#article-author-name")
		createdStr := e.ChildText("#article-create-date")
		createdStr = strings.TrimPrefix(createdStr, "Добавлен: ")

		createdTime, err := time.Parse("01-Oct-2022 08:15", createdStr)
		if err != nil {
			log.Fatal("time parse issue", err)
		}
		a.CreatedAt = createdTime
		a.Category = cat
		a.URL = "https://sber-invest.kz" + e.ChildAttr("a#article-title", "href")

		if a.Author != "" {
			ax = append(ax, a)
		}
	})

	url := fmt.Sprintf("https://sber-invest.kz/knowledgebase/%d", catN)

	c.Visit(url)

	return ax
}
