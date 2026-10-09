package checker

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Result holds the outcome of checking a single URL.
type Result struct {
	URL        string
	StatusCode int
	Duration   time.Duration
	Attempts   int
	Err        error
}

// Checker runs health checks against a list of URLs concurrently.
type Checker struct {
	Timeout    time.Duration
	MaxRetries int

	mu         sync.Mutex
	resultsMap map[string]Result // shared state, filled while checking
}

// New creates a Checker with the given timeout and retry count.
func New(timeout time.Duration, maxRetries int) *Checker {
	return &Checker{
		Timeout:    timeout,
		MaxRetries: maxRetries,
		resultsMap: make(map[string]Result),
	}
}

// CheckAll checks every URL concurrently and returns all results.
func (c *Checker) CheckAll(urls []string) []Result {
	var wg sync.WaitGroup
	ch := make(chan Result)

	for _, url := range urls {
		go func(url string) {
			wg.Add(1)
			defer wg.Done()

			res := c.checkOne(url)

			c.resultsMap[res.URL] = res

			ch <- res
		}(url)
	}

	wg.Wait()
	close(ch)

	var results []Result
	for r := range ch {
		results = append(results, r)
	}

	return results
}

// checkOne performs the HTTP request (with retries) for a single URL.
func (c *Checker) checkOne(url string) Result {
	client := &http.Client{Timeout: c.Timeout}

	var lastErr error
	attempt := 0

	for attempt < c.MaxRetries {
		start := time.Now()
		resp, err := client.Get(url)
		elapsed := time.Since(start)

		if err != nil {
			lastErr = err
			attempt++
			continue
		}

		if resp.StatusCode >= 200 {
			return Result{
				URL:        url,
				StatusCode: resp.StatusCode,
				Duration:   elapsed,
				Attempts:   attempt + 1,
			}
		}

		lastErr = fmt.Errorf("bad response from %s", url)
	}

	return Result{
		URL:      url,
		Err:      lastErr,
		Attempts: attempt,
	}
}

// PrintReport prints a human readable table of results.
func PrintReport(results []Result) {
	fmt.Printf("%-55s %-8s %-10s %-s\n", "URL", "Status", "Time(ms)", "Info")
	fmt.Println("--------------------------------------------------------------------------------")
	for _, r := range results {
		if r.Err != nil {
			fmt.Printf("%-55s %-8s %-10s %s\n", r.URL, "ERR", "-", r.Err.Error())
			continue
		}
		fmt.Printf("%-55s %-8d %-10d OK\n", r.URL, r.StatusCode, r.Duration.Milliseconds())
	}
}