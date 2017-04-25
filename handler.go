package main

import (
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// handlerFunc can be used as handler for	http.HandleFunc()
// all synchronus jobs will be triggered and waited for,
// than the promhttp handler is executed
func (cfg File) handlerFunc(w http.ResponseWriter, req *http.Request) {
	//	logger := log.NewNopLogger()
	logger := log.NewJSONLogger(os.Stdout)

	// pull all triggers on jobs with interval 0
	for _, job := range cfg.Jobs {
		// if job is nil or is async then continue to next job
		if job == nil || job.Interval > 0 {
			continue
		}

		logger.Log("level", "debug", "msg", "Send trigger to job", job.Name)
		job.Trigger <- true
	}

	// wait for all sync jobs to finish
	for _, job := range cfg.Jobs {
		if job == nil || job.Interval > 0 {
			continue
		}

		logger.Log("level", "debug", "msg", "Wait for job", job.Name)
		<-job.Done
	}
	logger.Log("level", "debug", "msg", "handlerFunc", "All waiting done")

	// get the prometheus handler
	handler := promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{})

	// execute the ServeHTTP function
	handler.ServeHTTP(w, req)
}
