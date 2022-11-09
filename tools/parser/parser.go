package parser

import (
	"fmt"

	"github.com/gocolly/colly"
)

func NewParser() {
	c := colly.NewCollector(
		colly.AllowedDomains("sber-invest.kz"),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnHTML("title", func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
	})

	c.Visit("https://sber-invest.kz/knowledgebase")

}
