package runner

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)


func RunBenchmark(cfg RunConfig) {
	fmt.Println("Starting benchmark...")

	var bodyData []byte
	var err error
	if cfg.BodyPath != "" {
		bodyData, err = os.ReadFile(cfg.BodyPath)
		if err != nil {
			fmt.Printf("Failed to read body file: %v\n", err)
			return
		}
	}

	client := &http.Client{
		Timeout: cfg.Timeout,
		Transport: &http.Transport{
			DisableKeepAlives: !cfg.KeepAlive,
		},
	}

	headerMap := parseHeaders(cfg.Headers)
	metrics := &Metrics{}
	var wg sync.WaitGroup
	var mu sync.Mutex

	requestsCh := make(chan struct{}, cfg.Concurrency*2) // buffered

	for range cfg.Concurrency {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range requestsCh {
				start := time.Now()
				req, err := http.NewRequest(cfg.Method, cfg.URL, bytes.NewReader(bodyData))
				if err != nil {
					mu.Lock()
					metrics.ErrorCount++
					metrics.TotalRequests++
					mu.Unlock()
					continue
				}
				for k, v := range headerMap {
					req.Header.Set(k, v)
				}

				resp, err := client.Do(req)
				latency := time.Since(start)

				mu.Lock()
				metrics.TotalRequests++
				metrics.Latencies = append(metrics.Latencies, latency)
				if err != nil {
					metrics.ErrorCount++
				} else {
					io.Copy(io.Discard, resp.Body)
					resp.Body.Close()
					if resp.StatusCode >= 200 && resp.StatusCode < 300 {
						metrics.SuccessCount++
					} else {
						metrics.ErrorCount++
					}
				}
				mu.Unlock()
			}
		}()
	}

	go func() {
		deadline := time.Now().Add(cfg.Duration)
		for time.Now().Before(deadline) {
			requestsCh <- struct{}{}
		}
		close(requestsCh)
	}()

	startTime := time.Now()
	wg.Wait()
	duration := time.Since(startTime)

	printSummary(metrics, duration)
}
