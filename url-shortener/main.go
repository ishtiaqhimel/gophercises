package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	fileName := flag.String("filename", "mapUrls.yaml", "YAML file of URLs list with shorten paths")
	flag.Parse()
	yamlHandler, err := YAMLHandler(*fileName, mapHandler)
	if err != nil {
		panic(err)
	}

	// Build the DBHandler using the dbHandler as the fallback
	loadDataToDB()
	dbHandler, err := DBHandler(yamlHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	log.Fatal(http.ListenAndServe(":8080", dbHandler))
}

func loadDataToDB() {
	db, err := bolt.Open("url.db", 0644, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	pathsToUrls := map[string]string{
		"/learn-linux": "https://linuxjourney.com",
		"/prometheus":  "https://a-cup-of.coffee/blog/prometheus",
	}

	for k, v := range pathsToUrls {
		if err := db.Update(func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists([]byte("urls"))
			if err != nil {
				log.Fatal(err)
			}

			return bucket.Put([]byte(k), []byte(v))
		}); err != nil {
			log.Fatal(err)
		}
	}
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World!")
}
