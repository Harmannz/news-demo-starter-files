package main

import (
	"fmt"
	"github.com/harmannz/news-demo-starter-files/news"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var pageTemplate = template.Must(template.ParseFiles("index.html"))

func indexHandler(responseWriter http.ResponseWriter, request *http.Request) {
	pageTemplate.Execute(responseWriter, nil)
}

func searchHandler(newsapi *news.Client) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		parsedURL, err := url.Parse(request.URL.String())
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		params := parsedURL.Query()
		searchQuery := params.Get("q")
		page := params.Get("page")
		if page == "" {
			page = "1"
		}

		results, err := newsapi.FetchEverything(searchQuery, page)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Printf("%+v", results)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == "" {
		log.Fatal("Env: apiKey must be set")
	}

	myClient := &http.Client{Timeout: 10 * time.Second}
	newsAPI := news.NewClient(myClient, apiKey, 20)

	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/search", searchHandler(newsAPI))

	fs := http.FileServer(http.Dir("assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.ListenAndServe(":"+port, mux)
}
