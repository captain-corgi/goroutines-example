package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/captain-corgi/goroutines-example/pkg/news"
)

const (
	API_KEY = "f32ee7b348msh230c75aaf106721p1366a6jsn952b266f7ae5"

	API_GOOGLE_NEWS_HOST = "google-news.p.rapidapi.com"
	GOOGLE_NEWS_URL      = "https://google-news.p.rapidapi.com/v1/top_headlines?lang=en&country=US"

	API_FREE_NEWS_HOST = "free-news.p.rapidapi.com"
	FREE_NEWS_URL      = "https://free-news.p.rapidapi.com/v1/search?lang=en&q=Elon"
)

func main() {
	// Greeting
	fmt.Println("Master goroutine and channel within 1 week")
	fmt.Println("This is an example of using select in Golang. It will wait for first response function")
	fmt.Println("We will have multiple function which return news")

	// Channels
	var (
		google = make(chan news.News)
		free   = make(chan news.News)
	)

	// Function list
	executors := []*news.Executor{
		{Fn: googleNews, Ch: google},
		{Fn: freeNews, Ch: free},
	}

	// Run each functions
	for _, executor := range executors {
		executor.Run()
	}

	// Wait for fastest response
	var articles []*news.Article
	select {
	case googleNewsResponse := <-google:
		fmt.Printf("Source: %s\n", googleNewsResponse.Source)
		articles = googleNewsResponse.Articles
	case freeNewsReponse := <-free:
		fmt.Printf("Source: %s\n", freeNewsReponse.Source)
		articles = freeNewsReponse.Articles
	}

	// Print result
	for i, arcticle := range articles {
		strByte, _ := json.MarshalIndent(arcticle, "", "  ")
		fmt.Printf("Articles %d:\n%s\n", i+1, strByte)
	}
}

func googleNews(google chan<- news.News) {
	req, err := http.NewRequest("GET", GOOGLE_NEWS_URL, nil)
	if err != nil {
		fmt.Printf("Error initializing request%v\n", err.Error())
		return
	}

	req.Header.Add("X-RapidAPI-Key", API_KEY)
	req.Header.Add("X-RapidAPI-Host", API_GOOGLE_NEWS_HOST)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request %v\n", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("Google News Response StatusCode %v Status %v\n", resp.StatusCode, resp.Status)
		return
	}

	googleNewsArticles := news.News{Source: "GoogleNewsApi"}
	if err := json.NewDecoder(resp.Body).Decode(&googleNewsArticles); err != nil {
		fmt.Printf("Error decoding body %v\n", err.Error())
		return
	}

	fmt.Printf("Google Articles Size %d\n", len(googleNewsArticles.Articles))
	google <- googleNewsArticles
}

func freeNews(free chan<- news.News) {
	req, err := http.NewRequest("GET", FREE_NEWS_URL, nil)
	if err != nil {
		fmt.Printf("Error initializing request%v\n", err.Error())
		return
	}

	req.Header.Add("X-RapidAPI-Key", API_KEY)
	req.Header.Add("X-RapidAPI-Host", API_FREE_NEWS_HOST)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request %v\n", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("Free News Response StatusCode %v Status %v\n", resp.StatusCode, resp.Status)
		return
	}

	var freeNewsArticles news.News
	if err := json.NewDecoder(resp.Body).Decode(&freeNewsArticles); err != nil {
		fmt.Printf("Error decoding body %v\n", err.Error())
		return
	}

	freeNewsArticles.Source = "FreeNewsApi"
	fmt.Printf("Free Articles Size %d\n", len(freeNewsArticles.Articles))
	free <- freeNewsArticles
}
