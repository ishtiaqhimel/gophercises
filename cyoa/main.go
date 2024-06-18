package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Story map[string]StoryArc

type StoryArc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text string `json:"text,omitempty"`
	Arc  string `json:"arc,omitempty"`
}

type Generate struct {
	story Story
	tpl   *template.Template
}

func main() {
	fileName := flag.String("file", "gopher.json", "the JSON file with the CYOA story")
	flag.Parse()

	data, err := os.ReadFile(*fileName)
	if err != nil {
		panic(err)
	}

	var story Story
	if err := json.Unmarshal(data, &story); err != nil {
		panic(err)
	}

	fmt.Println("Press y key to start the console and ENTER otherwise webserver will be started in 5 seconds...")
	keyEntered := make(chan string)
	timeOut := time.NewTimer(5 * time.Second)
	go readInput(keyEntered)

	select {
	case <-timeOut.C:
		generateForWeb(story)
	case input := <-keyEntered:
		if strings.HasPrefix(strings.ToLower(input), "y") {
			generateForConsole(story)
		} else {
			generateForWeb(story)
		}
	}
}

func readInput(c chan string) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	c <- input
}

func generateForWeb(story Story) {
	var gen Generate
	gen.story = story
	t := template.New("arcWeb.tpl")
	var err error
	gen.tpl, err = t.ParseFiles("arcWeb.tpl")
	if err != nil {
		panic(err)
	}
	gen.displayInWeb("intro")
}

func generateForConsole(story Story) {
	var gen Generate
	gen.story = story
	t := template.New("arc.tpl")
	var err error
	gen.tpl, err = t.ParseFiles("arc.tpl")
	if err != nil {
		panic(err)
	}
	gen.displayInConsole("intro")
}

func (g Generate) displayInConsole(arcName string) {
	storyArc, err := g.writeTemplatedText(os.Stdout, arcName)
	if err != nil {
		panic(err)
	}

	if len(storyArc.Options) == 0 {
		return
	}

	fmt.Print("Your Option: ")
	var optionNumber int
	fmt.Scan(&optionNumber)
	for id, option := range storyArc.Options {
		if id == optionNumber {
			g.displayInConsole(option.Arc)
		}
	}
}

func (g Generate) displayInWeb(arcName string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		arcs := r.URL.Query()["arc"]
		if len(arcs) > 0 {
			arcName = arcs[0]
		}

		_, err := g.writeTemplatedText(w, arcName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("Listening on 8080 port .....")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (g Generate) writeTemplatedText(w io.Writer, arcName string) (*StoryArc, error) {
	storyArc, err := g.getArc(arcName)
	if err != nil {
		return nil, err
	}

	if err = g.tpl.Execute(w, storyArc); err != nil {
		return nil, err
	}
	return storyArc, nil
}

func (g Generate) getArc(key string) (*StoryArc, error) {
	for k, v := range g.story {
		if k == key {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("cannot find %s arc", key)
}
