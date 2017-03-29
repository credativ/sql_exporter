package main

import (
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// handlerFunc can be used as handler for	http.HandleFunc()
// all synchronus jobs will be triggered and waited for,
// than the promhttp handler is executed
func handlerFunc(w http.ResponseWriter, req *http.Request) {
	//logger := log.NewJSONLogger(os.Stdout)
	logger := log.NewNopLogger()

	if httpTrigger != nil && httpDone != nil {
		// pull the trigger and submit the run count
		logger.Log("level", "debug", "msg", "pull the trigger")
		httpTrigger <- true

		// wait for synchronus actions to finish
		logger.Log("level", "debug", "msg", "wait for done")
		<-httpDone
	}

	// get the prometheus handler
	handler := promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{})

	// execute die ServeHTTP function
	handler.ServeHTTP(w, req)
}

// syncHandler handles the global trigger to execute all synchronus tasks
func syncHandler(trigger chan bool, done chan bool, cfg *File) {
	//logger := log.NewJSONLogger(os.Stdout)
	logger := log.NewNopLogger()

	for {
		// wait for trigger
		<-trigger

		// pull all triggers on jobs with interval 0
		for _, job := range cfg.Jobs {
			if job == nil {
				continue
			}

			if job.Interval <= 0 {
				logger.Log("level", "debug", "msg", "Send trigger to job", job.Name)
				job.Trigger <- true
			}
		}

		// wait for all jobs to finish
		for _, job := range cfg.Jobs {
			if job == nil {
				continue
			}

			if job.Interval <= 0 {
				logger.Log("level", "debug", "msg", "Wait for job", job.Name)
				<-job.Done
			}
		}
		logger.Log("level", "debug", "msg", "syncHandler", "All waiting done")

		// send global done so the handler will continue
		done <- true
	}
}
