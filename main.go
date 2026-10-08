package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"urlchecker/checker"
)

func main() {
	urls, err := readURLs("urls.txt")
	if err != nil {
		fmt.Println("failed to read urls.txt:", err)
		os.Exit(1)
	}

	fmt.Printf("Checking %d URLs...\n\n", len(urls))

	c := checker.New(5*time.Second, 3)

	start := time.Now()
	results := c.CheckAll(urls)
	totalTime := time.Since(start)

	checker.PrintReport(results)

	fmt.Printf("\nChecked %d URLs in %v\n", len(results), totalTime)
}

func readURLs(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var urls []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		urls = append(urls, line)
	}

	return urls, scanner.Err()
}
