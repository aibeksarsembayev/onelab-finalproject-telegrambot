package parser

import (
	"fmt"

	"github.com/gocolly/colly"
)

type Article struct {
	Title string
	// Content   string
	Author    string
	Category  string
	URL       string
	CreatedAt string
}

func NewParser() {
	c := colly.NewCollector(
		colly.AllowedDomains("sber-invest.kz"),
	)

	ax := []*Article{}

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnHTML("div.message", func(e *colly.HTMLElement) {
		a := &Article{}

		a.Title = e.ChildText("#article-title")
		a.Author = e.ChildText("#article-author-name")
		a.CreatedAt = e.ChildText("#article-create-date")
		a.Category = "category name text"
		a.URL = "https://sber-invest.kz" + e.ChildAttr("a#article-title", "href")
		ax = append(ax, a)
	})

	c.Visit("https://sber-invest.kz/knowledgebase/1")

	for i, a := range ax {
		fmt.Println("article #", i, a)
	}

}
