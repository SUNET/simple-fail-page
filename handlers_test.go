package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"io/ioutil"
)

func TestServeFileHandlerIndex(t *testing.T) {
	configData, err := ioutil.ReadFile("simple-fail-page.yaml")
	check(err)
	configuration := readConfig(configData)

	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response
	rr := httptest.NewRecorder()

	// Setup handler
	handler := http.Handler(serveFile(configuration))

	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/", nil)
	check(err)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><title>Fail</title><link rel='stylesheet' id='main-css' href='assets/style.css' type='text/css' media='all' /></head><body><h1>Panic</h1></body></html>`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestServeFileHandlerAssets(t *testing.T) {
	configData, err := ioutil.ReadFile("simple-fail-page.yaml")
	check(err)
	configuration := readConfig(configData)

	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response
	rr := httptest.NewRecorder()

	// Setup handler
	handler := http.Handler(serveFile(configuration))

	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/assets/style.css", nil)
	check(err)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `body { background-color: red; }`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestServeFileHandlerNotFound(t *testing.T) {
	configuration := Configuration{}
	configuration.Redirect404 = false

	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response
	rr := httptest.NewRecorder()

	// Setup handler
	handler := http.Handler(serveFile(configuration))

	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/assets/style.css", nil)
	check(err)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestServeFileHandlerRedirect(t *testing.T) {
	configuration := Configuration{}
	configuration.Redirect404 = true

	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response
	rr := httptest.NewRecorder()

	// Setup handler
	handler := http.Handler(serveFile(configuration))

	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/assets/style.css", nil)
	check(err)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusMovedPermanently {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestServeFileHandlerJSON(t *testing.T) {
	configuration := Configuration{JsonResponse: map[string]string{"test": "test"}}

	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response
	rr := httptest.NewRecorder()

	// Setup handler
	handler := http.Handler(serveFile(configuration))

	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/assets/style.css", nil)
	check(err)
	req.Header.Set("Content-Type", "application/json")

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}