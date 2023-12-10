package main

import (
	"awesomeProject1/internal/handlers"
	"github.com/gin-gonic/gin"
)

var globalMemory = handlers.MemStorage{}

func main() {
	globalMemory.Counter = make(map[string]int64)
	globalMemory.Gauge = make(map[string]float64)

	r := gin.Default()
	r.POST("/update/:type/:name/:value", func(c *gin.Context) {
		globalMemory.UpdateMetrics(c)
	})
	r.GET("/value/:type/:name/", func(c *gin.Context) {
		globalMemory.GetMetric(c)
	})
	r.GET("/", func(c *gin.Context) {
		globalMemory.PrintMetrics(c)
	})

	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
