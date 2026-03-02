package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	var urls []string
	urls = HandleInput()

	// Sequential
	fmt.Println("Sequential fetching")
	seqStart := time.Now()
	results := CheckHealth(urls)
	for _, result := range results {
		fmt.Println(result)
	}
	seqElapsed := time.Since(seqStart)

	// Concurrent
	fmt.Println("Concurrency fetching")
	jobs := make(chan string, 5)
	resultsChan := make(chan string, 5)

	concStart := time.Now()

	for w := 0; w < 5; w++ {
		go Worker(w, jobs, resultsChan)
	}

	for _, url := range urls {
		jobs <- url
	}
	close(jobs)

	for range urls {
		fmt.Println(<-resultsChan)
	}
	close(resultsChan)
	concElapsed := time.Since(concStart)

	fmt.Println("Sequential:", seqElapsed, "| Concurrent:", concElapsed)
}

var client = &http.Client{Timeout: 5 * time.Second}

func HandleInput() []string {
	var urls []string

	Scanner := bufio.NewScanner(os.Stdin)
	var inputUrls string
	if Scanner.Scan() {
		inputUrls = Scanner.Text()
	}

	var url string
	for _, c := range inputUrls {
		if c == ' ' {
			urls = append(urls, url)
			url = ""
			continue
		}
		url = url + string(c)
	}

	urls = append(urls, url)

	return urls
}

func CheckHealth(urls []string) []string {
	var results []string
	for _, url := range urls {
		res, err := client.Get(url)
		var result string
		if err != nil {
			result = "Failed to fetch"
		} else {
			defer res.Body.Close()
			result = res.Status
		}
		results = append(results, result)
	}

	return results
}

func Worker(id int, jobs <-chan string, results chan<- string) {
	for job := range jobs {
		res, err := client.Get(job)
		fmt.Println("Worker ", id, "started job:", job)
		if err != nil {
			results <- fmt.Sprintf("%s -> ERROR: %v", job, err)
		} else {
			defer res.Body.Close()
			results <- res.Status
		}
	}
}
