package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robfig/cron/v3"
)

// Keep counters and status here. (TODO: persist)
type Response struct {
	Schedule1  int
	Schedule2  int
	ScratchJob int
}

var response = &Response{}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

var (
	uptime = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sumo_test_metric_admin_uptime",
		Help: "The uptime of Sumo Admin (scheduled jobs server)",
	})
)

var (
	scratchJob = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "sumo_test_gauge_admin_scratch_job",
		Help: "The result of last scratch delete job (scheduled jobs server)",
	})
)

func main() {
	c := cron.New()

	// Schedule 1: Every other second, counter
	c.Schedule(cron.Every(2*time.Second), cron.FuncJob(func() {
		uptime.Inc()
		response.Schedule1++
	}))

	// Schedule 2: Every five minute, counter (https://crontab.guru/)
	c.AddFunc("*/5 * * * *", func() {
		response.Schedule2++
	})

	// Schedule 3: The scratch delete job
	// Return negative on error, positive on successful run
	c.AddFunc("*/5 * * * *", func() {
		// Run the scratch delete job... TODO
		if time.Now().Minute() < 30 {
			scratchJob.Set(-2)
			response.ScratchJob = -2
		} else {
			scratchJob.Set(2)
			response.ScratchJob = 2
		}
	})

	c.Start()

	// Turn web server on to show values on web page too
	r := gin.New()
	r.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, response)
	})

	// Prometheus to scrape our metrics here
	r.GET("/metrics", prometheusHandler())

	r.Run(":8000")
}
