package main

import (
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
)

// Build time vars
var (
	Name        = "prom-sql-exporter"
	Version     string
	BuildTime   string
	Commit      string
	cfg         File
	httpTrigger chan bool // Global trigger for sync http actions
	httpDone    chan bool // Global done to signal back to http handler

)

func main() {
	// init logger
	logger := log.NewJSONLogger(os.Stdout)
	logger = log.NewContext(logger).With(
		"ts", log.DefaultTimestampUTC,
		"caller", log.DefaultCaller,
		"name", Name,
		"version", Version,
		"commit", Commit,
	)

	cfgFile := "config.yml"
	if f := os.Getenv("CONFIG"); f != "" {
		cfgFile = f
	}

	// read config
	cfg, err := Read(cfgFile)
	if err != nil {
		panic(err)
	}

	// dispatch all jobs
	for _, job := range cfg.Jobs {
		if job == nil {
			continue
		}
		job.log = log.NewContext(logger).With("job", job.Name)
		go job.Run()
	}

	// create needed channels
	httpTrigger = make(chan bool)
	httpDone = make(chan bool)

	// start handler for synchonus http actions
	go syncHandler(httpTrigger, httpDone, &cfg)

	http.HandleFunc("/metrics", handlerFunc)
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { http.Error(w, "OK", http.StatusOK) })
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
		<head><title>SQL Exporter</title></head>
		<body>
		<h1>SQL Exporter</h1>
		<p><a href="/metrics">Metrics</a></p>
		</body>
		</html>
		`))
	})

	addr := ":8080"
	logger.Log("level", "info", "msg", "Starting sql_exporter", "addr", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		logger.Log("level", "error", "msg", "Error starting HTTP server:", "err", err)
		os.Exit(1)
	}
}
