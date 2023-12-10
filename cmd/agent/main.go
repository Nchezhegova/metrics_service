package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"
)

// Функция для отправки метрики на сервер
func sendMetric(metricType, name string, value interface{}) {
	url := fmt.Sprintf("http://localhost:8080/update/%s/%s/%v", metricType, name, value)
	_, err := http.Post(url, "text/plain", nil)
	if err != nil {
		fmt.Println("Error sending metric:", err)
	}
}

func collectMetrics(metrics *map[string]interface{}, mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	(*metrics)["Alloc"] = memStats.Alloc
	(*metrics)["BuckHashSys"] = memStats.BuckHashSys
	(*metrics)["Frees"] = memStats.Frees
	(*metrics)["GCCPUFraction"] = memStats.GCCPUFraction
	(*metrics)["GCSys"] = memStats.GCSys
	(*metrics)["HeapAlloc"] = memStats.HeapAlloc
	(*metrics)["HeapIdle"] = memStats.HeapIdle
	(*metrics)["HeapInuse"] = memStats.HeapInuse
	(*metrics)["HeapObjects"] = memStats.HeapObjects
	(*metrics)["HeapReleased"] = memStats.HeapReleased
	(*metrics)["HeapSys"] = memStats.HeapSys
	(*metrics)["LastGC"] = memStats.LastGC
	(*metrics)["Lookups"] = memStats.Lookups
	(*metrics)["MCacheInuse"] = memStats.MCacheInuse
	(*metrics)["MCacheSys"] = memStats.MCacheSys
	(*metrics)["MSpanInuse"] = memStats.MSpanInuse
	(*metrics)["MSpanSys"] = memStats.MSpanSys
	(*metrics)["Mallocs"] = memStats.Mallocs
	(*metrics)["NextGC"] = memStats.NextGC
	(*metrics)["NumForcedGC"] = memStats.NumForcedGC
	(*metrics)["NumGC"] = memStats.NumGC
	(*metrics)["OtherSys"] = memStats.OtherSys
	(*metrics)["PauseTotalNs"] = memStats.PauseTotalNs
	(*metrics)["StackInuse"] = memStats.StackInuse
	(*metrics)["StackSys"] = memStats.StackSys
	(*metrics)["Sys"] = memStats.Sys
	(*metrics)["TotalAlloc"] = memStats.TotalAlloc
	// ещё 2 метрики
	(*metrics)["RandomValue"] = rand.Float64()

}

func main() {
	pollInterval := 2 * time.Second
	reportInterval := 10 * time.Second

	var pollCount int64
	metrics := make(map[string]interface{})
	var mu sync.Mutex

	// Таймер для сбора метрик
	go func() {
		for {
			collectMetrics(&metrics, &mu)
			pollCount++
			time.Sleep(pollInterval)
		}
	}()

	// Таймер для отправки метрик
	for {
		<-time.After(reportInterval)

		mu.Lock()
		for name, value := range metrics {
			go sendMetric("gauge", name, value)
		}
		sendMetric("counter", "PollCount", pollCount)
		mu.Unlock()
	}
}
