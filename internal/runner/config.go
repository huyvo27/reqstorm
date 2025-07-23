package runner

import "time"

type RunConfig struct {
	URL         string
	Method      string
	Concurrency int
	Duration    time.Duration
	Headers     []string
	BodyPath    string
	Timeout     time.Duration
	ReportType  string
	KeepAlive   bool
}

type Metrics struct {
	TotalRequests int
	SuccessCount  int
	ErrorCount    int
	Latencies     []time.Duration
}