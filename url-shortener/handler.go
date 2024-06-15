package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/go-yaml/yaml"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(fileName string, fallback http.Handler) (http.HandlerFunc, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open url yaml file: %w", err)
	}
	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read url yaml file: %w", err)
	}

	var pathUrls []pathUrl
	if err = yaml.Unmarshal(bytes, &pathUrls); err != nil {
		return nil, fmt.Errorf("failed to unmarhal: %w", err)
	}

	pathsToUrls := make(map[string]string)
	for _, pu := range pathUrls {
		pathsToUrls[pu.Path] = pu.Url
	}
	return MapHandler(pathsToUrls, fallback), nil
}

type pathUrl struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

func DBHandler(fallback http.Handler) (http.HandlerFunc, error) {
	db, err := bolt.Open("url.db", 0644, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	defer db.Close()

	pathsToUrls := make(map[string]string)
	if err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("urls"))
		if bucket == nil {
			return fmt.Errorf("bucket urls not found")
		}
		if err := bucket.ForEach(func(k, v []byte) error {
			pathsToUrls[string(k)] = string(v)
			return nil
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return MapHandler(pathsToUrls, fallback), nil
}
