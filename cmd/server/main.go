package main

import (
	"fmt"
	"net/http"
	"sync"
)

type MemStorage struct {
	gauge   map[string]float64 `json:"gauge"`
	counter map[string]int64   `json:"counter"`
}

func (m *MemStorage) countValue(k string, v int64) {
	mu.Lock()
	defer mu.Unlock()
	m.counter[k] += v
	fmt.Println(m.counter)
}
func (m *MemStorage) gaugeValue(k string, v float64) {
	mu.Lock()
	defer mu.Unlock()
	m.gauge[k] = v
	fmt.Println(m.gauge)
}

var globalMemory = MemStorage{}
var mu sync.Mutex

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/counter/`, counterMetric)
	mux.HandleFunc(`/update/gauge/`, gaugeMetric)

	globalMemory.counter = make(map[string]int64)
	globalMemory.gauge = make(map[string]float64)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
