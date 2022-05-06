package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

type Response struct {
	Schedule1 int
	Schedule2 int
}

var response = &Response{}

func main() {
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
	r.Run(":8000")
}
