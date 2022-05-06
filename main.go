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

type Response struct {
	Schedule1 int
	Schedule2 int
}

var response = &Response{}

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sumo_metric_test_daily_job_success",
		Help: "The success of daily job",
	})
)

func main() {
	recordMetrics()

	c := cron.New()

	// Schedule 1
	c.Schedule(cron.Every(2*time.Second), cron.FuncJob(func() {
		response.Schedule1++
	}))

	// Schedule 2: Every five minute - https://crontab.guru/
	c.AddFunc("*/5 * * * *", func() {
		response.Schedule2++
	})
	c.Start()

	// Tun web server to show values
	r := gin.New()
	r.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, response)
	})
	r.GET("/howdy", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "howdy")
	})
	r.GET("/metrics", func(ctx *gin.Context) {
		promhttp.Handler()
	})

	//http.Handle("/metrics", promhttp.Handler())

	r.Run(":8000")
}
