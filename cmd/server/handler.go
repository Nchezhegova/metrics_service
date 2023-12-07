package main

import (
	"net/http"
	"strconv"
	"strings"
)

func splitUrl(u string) []string {
	//можно было использовать гориллу, но вроде просили без сторонних библиотек
	parts := strings.Split(u, "/")
	return parts
}

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
	globalMemory.countValue(k, v)
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
	globalMemory.gaugeValue(k, v)
	res.WriteHeader(http.StatusOK)
}
