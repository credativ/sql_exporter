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
func (ex *Exporter) handlerFunc(w http.ResponseWriter, req *http.Request) {
	//	logger := log.NewNopLogger()
	logger := log.NewJSONLogger(os.Stdout)
	logger = log.With(logger, "caller", "handlerFunc")

	// pull all triggers on jobs with interval 0
	logger.Log("level", "debug", "msg", "Start all sync jobs")
	for _, job := range ex.jobs {
		// if job is nil or is async then continue to next job
		if job == nil || job.Interval > 0 {
			logger.Log("level", "debug", "msg", "Send NO trigger to job", job.Name)
			continue
		}
		logger.Log("level", "debug", "msg", "Send trigger to job", job.Name)
		job.Trigger <- true
	}

	// wait for all sync jobs to finish
	for _, job := range ex.jobs {
		if job == nil || job.Interval > 0 {
			continue
		}
		logger.Log("level", "debug", "msg", "Wait for job", job.Name)
		<-job.Done
	}
	logger.Log("level", "debug", "msg", "All waiting done")

	// get the prometheus handler
	handler := promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{})

	// execute the ServeHTTP function
	handler.ServeHTTP(w, req)
}
