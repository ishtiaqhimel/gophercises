package main

import (
	"flag"
	"fmt"
	"hlp/link"
	"os"
)

func main() {
	fileName := flag.String("file", "example.html", "The html file to parse links from")
	flag.Parse()
	f, err := os.Open(*fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	links, err := link.Parse(f)
	if err != nil {
		panic(err)
	}

	fmt.Println(links)
}
