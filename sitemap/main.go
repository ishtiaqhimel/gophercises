package main

import (
	"encoding/xml"
	"flag"
	"io"
	"net/http"
	"net/url"
	"os"
	"sitemap/link"
	"strings"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlSet struct {
	Xmlns string `xml:"xmlns,attr"`
	Urls  []loc  `xml:"url"`
}

func main() {
	urlStr := flag.String("url", "https://gophercises.com", "the url that you want to build a sitemap for")
	maxDepth := flag.Int("depth", 10, "the maximum number of links deep to traverse")
	flag.Parse()

	pages := bfs(*urlStr, *maxDepth)
	toXml := urlSet{
		Xmlns: xmlns,
	}
	for _, page := range pages {
		toXml.Urls = append(toXml.Urls, loc{page})
	}

	out, err := xml.MarshalIndent(toXml, "", "  ")
	if err != nil {
		panic(err)
	}

	// Printing XML
	os.Stdout.Write([]byte(xml.Header))
	os.Stdout.Write(out)
	os.Stdout.Write([]byte("\n"))
}

func bfs(urlStr string, maxDepth int) []string {
	visited := make(map[string]struct{})
	var urls []string
	curr := []string{urlStr}
	for i := 0; i <= maxDepth; i++ {
		urls, curr = curr, make([]string, 0)
		if len(urls) == 0 {
			break
		}
		for _, url := range urls {
			if _, ok := visited[url]; ok {
				continue
			}
			visited[url] = struct{}{}
			curr = append(curr, get(url)...)
		}
	}
	pages := make([]string, 0, len(visited))
	for url := range visited {
		pages = append(pages, url)
	}
	return pages
}

func get(urlStr string) []string {
	res, err := http.Get(urlStr)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	reqUrl := res.Request.URL
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	base := baseUrl.String()

	return filter(getHrefs(res.Body, base), withBasePrefix(base))
}

func getHrefs(r io.Reader, base string) []string {
	links, _ := link.Parse(r)
	var hrefs []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			hrefs = append(hrefs, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			hrefs = append(hrefs, l.Href)
		}
	}
	return hrefs
}

func filter(links []string, fn func(link string) bool) []string {
	var ret []string
	for _, link := range links {
		if fn(link) {
			ret = append(ret, strings.TrimSuffix(link, "/"))
		}
	}

	return ret
}

func withBasePrefix(base string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, base)
	}
}
