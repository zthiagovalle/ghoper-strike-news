package feed

import (
	"time"

	"github.com/mmcdole/gofeed"
)

type Item struct {
	Title           string
	Link            string
	Description     string
	Content         string
	PublishedParsed *time.Time
}

type RSSParser struct {
	parser *gofeed.Parser
	url    string
}

func NewParser(url string) *RSSParser {
	return &RSSParser{
		parser: gofeed.NewParser(),
		url:    url,
	}
}

func (p *RSSParser) Parse() ([]*Item, error) {
	feed, err := p.parser.ParseURL(p.url)
	if err != nil {
		return nil, err
	}

	limit := 3
	if len(feed.Items) < limit {
		limit = len(feed.Items)
	}

	items := make([]*Item, 0, limit)
	for i := 0; i < limit; i++ {
		fi := feed.Items[i]
		items = append(items, &Item{
			Title:           fi.Title,
			Link:            fi.Link,
			Description:     fi.Description,
			Content:         fi.Content,
			PublishedParsed: fi.PublishedParsed,
		})
	}

	return items, nil
}
