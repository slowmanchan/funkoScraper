package scraper

import (
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

type Collection struct {
	Funkos []*Funko
}

type Funko struct {
	Name     string
	ImgURL   string
	Brand    string
	Series   string
	Produced string
	Scale    string
	Edition  string
}

func Run() {
	collection := new(Collection)
	pageCount := 1
	pageStr := ""

	for {
		c := colly.NewCollector()

		c.OnHTML(".catalog-item-search-results div", func(e *colly.HTMLElement) {
			if e.Text == "" {
				fmt.Println("finished scraping pages")
				os.Exit(1)
			}
			collection.Funkos = append(collection.Funkos, NewFunko(e))
		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL)
			pageCount++
		})

		if pageCount > 1 {
			pageStr = fmt.Sprintf(`&page=%d`, pageCount)
		}

		c.Visit("https://www.hobbydb.com/catalog_item_types/vinyl-art-toys/keyword_type_search?keyword_query=funko+pop" + pageStr)
	}
}

func NewFunko(e *colly.HTMLElement) *Funko {
	funko := new(Funko)
	e.ForEach(".search-result-field-list li", func(i int, e *colly.HTMLElement) {
		str := strings.Split(e.Text, ":")

		switch str[0] {
		case "Brand":
			funko.Brand = str[1]
		case "Series":
			funko.Series = str[1]
		case "Produced":
			funko.Produced = strings.Join(str[1:], " ")
		case "Scale":
			funko.Scale = str[1]
		}
	})

	funko.Name = e.ChildAttr(".image-container a img", "alt")
	funko.ImgURL = e.ChildAttr(".image-container a img", "src")

	return funko
}
