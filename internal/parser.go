package internal

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/antchfx/xmlquery"
)

// Parser represents
type Parser struct {
	logger *log.Logger
}

// NewParser
func NewParser() *Parser {
	return &Parser{
		log.New(os.Stderr, parserLogPrefix, loggerFlags),
	}
}

// Parse
func (p *Parser) Parse(ctx context.Context, feeds <-chan Feed, posts chan<- Post) {
	for feed := range feeds {
		p.logger.Printf("looking for posts in feed %s", feed.name)

		doc, err := xmlquery.Parse(feed.content)
		if err != nil {
			continue
		}

		items, err := xmlquery.QueryAll(doc, feed.config.Rules.PostPath)
		if err != nil {
			continue
		}

		for _, item := range items {
			post := Post{}

			post.Title, err = getElementData(item, feed.config.Rules.TitlePath)
			if err != nil {
				p.logger.Printf("could not parse post title: %s", err)
			}

			post.Link, err = getElementData(item, feed.config.Rules.LinkPath)
			if err != nil {
				p.logger.Printf("could not parse post link: %s", err)
			}

			post.Description, err = getElementData(item, feed.config.Rules.DescriptionPath)
			if err != nil {
				p.logger.Printf("could not parse post description")
			}

			posts <- post
		}
	}

	log.Println("closing posts channel")
	close(posts)
}

// getElementData
func getElementData(el *xmlquery.Node, path string) (string, error) {
	child := el.SelectElement(path)
	if child == nil {
		return "", errors.New("element not found")
	}

	if !strings.Contains(path, "@") {
		return child.InnerText(), nil
	}

	parts := strings.Split(path, "@")
	attr := parts[1]

	return child.SelectAttr(attr), nil
}

// findItems
// func (p *Parser) findItems(feed *Feed, posts chan<- Post) {
// 	decoder := xml.NewDecoder(feed.content)

// feed:
// 	for {
// 		token, err := decoder.Token()

// 		switch {
// 		case errors.Is(err, io.EOF):
// 			break feed

// 		case err != nil:
// 			continue
// 		}

// 		st, ok := token.(xml.StartElement)
// 		if ok && st.Name.Local == feed.config.Rules.PostTag {
// 			post, err := p.initPost(decoder, feed.config.Rules)
// 			if err != nil {
// 				continue feed
// 			}
// 			posts <- post
// 		}
// 	}
// }

// // initPost
// func (p *Parser) initPost(decoder *xml.Decoder, rules Rules) (Post, error) {
// 	var post Post
// 	for {
// 		token, err := decoder.Token()
// 		if err != nil {
// 			continue
// 		}

// 		st, ok := token.(xml.StartElement)
// 		if !ok {
// 			continue
// 		}

// 		switch st.Name.Local {
// 		case rules.TitleTag:
// 			t, err := decoder.Token()
// 			if err != nil {
// 				continue
// 			}

// 			cd, ok := t.(xml.CharData)
// 			if !ok {
// 				continue
// 			}

// 			post.Title = string(cd)
// 			p.logger.Printf("got title: %q", string(cd))
// 		case rules.LinkTag:
// 			p.logger.Printf("link attrs: %#v", st.Attr)
// 			// post.Link = string(cd)
// 			// p.logger.Printf("got link: %q", string(cd))
// 		}
// 		break
// 	}
// 	return post, nil
// }
