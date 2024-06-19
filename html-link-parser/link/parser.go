package link

import (
	"io"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	return dfs(doc), nil
}

func dfs(n *html.Node) []Link {
	var links []Link
	if n.Type == html.ElementNode &&
		n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, Link{a.Val, getText(n)})
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = append(links, dfs(c)...)
	}
	return links
}

func getText(n *html.Node) string {
	text := ""
	if n.Type == html.TextNode {
		return n.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += getText(c)
	}
	return text
}
