package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/huyvo27/reqstorm/internal/runner"
)

var (
	url         string
	concurrency int
	duration    time.Duration
	method      string
	headers     []string
	bodyPath    string
	timeout     time.Duration
	reportType  string
	keepAlive   bool
)

var rootCmd = &cobra.Command{
	Use:   "reqstorm",
	Short: "A simple CLI tool to benchmark HTTP endpoints",
	Run: func(cmd *cobra.Command, args []string) {
		if url == "" {
			fmt.Println("Error: --url is required")
			os.Exit(1)
		}

		runner.RunBenchmark(runner.RunConfig{
			URL:         url,
			Method:      method,
			Concurrency: concurrency,
			Duration:    duration,
			Headers:     headers,
			BodyPath:    bodyPath,
			Timeout:     timeout,
			ReportType:  reportType,
			KeepAlive:   keepAlive,
		})
	},
}

func init() {
	rootCmd.Flags().StringVarP(&url, "url", "u", "", "Target URL to benchmark (required)")
	rootCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 10, "Number of concurrent users")
	rootCmd.Flags().DurationVarP(&duration, "duration", "d", 10*time.Second, "Duration of the test (e.g. 10s, 1m)")
	rootCmd.Flags().StringVarP(&method, "method", "m", "GET", "HTTP method to use")
	rootCmd.Flags().StringSliceVarP(&headers, "header", "H", []string{}, "Custom headers (can be used multiple times)")
	rootCmd.Flags().StringVarP(&bodyPath, "body", "b", "", "Request body file path (JSON, etc)")
	rootCmd.Flags().DurationVar(&timeout, "timeout", 5*time.Second, "Request timeout duration")
	rootCmd.Flags().StringVar(&reportType, "report", "text", "Output format (text|json)")
	rootCmd.Flags().BoolVar(&keepAlive, "keep-alive", true, "Use HTTP keep-alive")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
