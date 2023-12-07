package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type MemStorage struct {
	gauge   map[string]float64 `json:"gauge"`
	counter map[string]int64   `json:"counter"`
}

type Storage interface {
	countValue()
	gaugeValue()
}

func (m MemStorage) countValue(k string, v int64) {
	m.counter[k] += v
	fmt.Println(m.counter)
}
func (m MemStorage) gaugeValue(k string, v float64) {
	m.gauge[k] = v
	fmt.Println(m.gauge)
}

func splitUrl(u string) []string {
	//можно было использовать гориллу, но вроде просили без сторонних библиотек
	parts := strings.Split(u, "/")
	fmt.Println(len(parts))
	return parts
}

var m = MemStorage{}

func counterMetric(res http.ResponseWriter, req *http.Request) {
	//новое значение добавляется к предыдущему
	if req.Method != "POST" {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	parts := splitUrl(req.URL.Path)

	if len(parts) != 5 {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	v, err := strconv.ParseInt(parts[len(parts)-1], 10, 64)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	k := parts[len(parts)-2]
	m.countValue(k, v)
	res.WriteHeader(http.StatusOK)
}
func gaugeMetric(res http.ResponseWriter, req *http.Request) {
	//новое значение замещает предыдущее, если было известно
	if req.Method != "POST" {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	parts := splitUrl(req.URL.Path)
	if len(parts) != 5 {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	v, err := strconv.ParseFloat(parts[len(parts)-1], 64)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	k := parts[len(parts)-2]
	m.gaugeValue(k, v)
	res.WriteHeader(http.StatusOK)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/counter/`, counterMetric)
	mux.HandleFunc(`/update/gauge/`, gaugeMetric)

	m.counter = make(map[string]int64)
	m.gauge = make(map[string]float64)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
