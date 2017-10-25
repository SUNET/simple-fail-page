package main

// Inspired by
// https://gist.github.com/superbrothers/0a8b6390c6315916aeb8

import (
	"flag"
	"log"
	"net/http"
	"path"
	"os"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var name string = "simple-fail-page"
var version string = "0.0.1"

type Configuration struct {
	UrlPathToFilePath map[string]string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readConfig(configData []byte) Configuration {
	config := Configuration{}
	err := yaml.Unmarshal(configData, &config)
	check(err)
	return config
}

func serveFile(urlPathToFilePath map[string]string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		filePath := urlPathToFilePath[r.URL.Path]
		if filePath != "" {
			fileStat, err := os.Stat(filePath)
			check(err)
			file, err := os.Open(filePath)
			check(err)
			defer file.Close()
			_, filename := path.Split(filePath)
			modTime := fileStat.ModTime()
			log.Println(r.URL.Path)
			http.ServeContent(w, r, filename, modTime, file)
		} else {
			log.Println("Access denied: " + r.URL.Path)
			http.NotFound(w, r)
		}
	}
	return http.HandlerFunc(fn)
}

func main() {
	log.Printf("%s %s", name, version)
	// Handle args
	var listenTo string
	var configFile string
	flag.StringVar(&listenTo, "listen", ":80", "[address]:port")
	flag.StringVar(&configFile, "config", "simple-fail-page.yaml", "Path to YAML config file")
	flag.Parse()
	// Read config
	configData, err := ioutil.ReadFile(configFile)
	check(err)
	configuration := readConfig(configData)
	// Set up handles for requests
	http.Handle("/", serveFile(configuration.UrlPathToFilePath))
	// Run server
	log.Fatal(http.ListenAndServe(listenTo, nil))
}
