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
	"strings"
	"encoding/json"
)

var name string = "simple-fail-page"
var version string = "0.0.1"

type Configuration struct {
	UrlPathToFilePath map[string]string
	Redirect404 bool
	JsonResponse map[string]string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func readConfig(configData []byte) Configuration {
	config := Configuration{}
	err := yaml.Unmarshal(configData, &config)
	check(err)
	return config
}

func checkRequestType(r *http.Request) string {
	defaultType := "*/*"
	jsonType := "application/json"
	contentType := strings.Split(r.Header.Get("Content-Type"), ",")
	accepts := strings.Split(r.Header.Get("Accept"), ",")
	if (stringInSlice(jsonType, contentType) || stringInSlice(jsonType, accepts)) {
		return jsonType
	} else {
		return defaultType
	}
}

func createJsonResponse(config Configuration, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	jsonResponseData, err := json.Marshal(config.JsonResponse)
	check(err)
	w.Write(jsonResponseData)
}

func serveFile(config Configuration) http.Handler {
	urlPathToFilePath := config.UrlPathToFilePath
	redirect404 := config.Redirect404
	fn := func(w http.ResponseWriter, r *http.Request) {
		if checkRequestType(r) == "application/json" && len(config.JsonResponse) != 0{
			createJsonResponse(config, w)
		} else {
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
				if redirect404 == true {
					log.Println("Redirecting: " + r.URL.Path + " -> /")
					http.Redirect(w, r, "/", http.StatusMovedPermanently)
				} else {
					http.NotFound(w, r)
				}
			}
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
	http.Handle("/", serveFile(configuration))
	// Run server
	log.Fatal(http.ListenAndServe(listenTo, nil))
}
