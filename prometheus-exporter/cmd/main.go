package main

import (
	"fmt"
	"net/http"

	"github.com/lizw91/monitor/prometheus-exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// 创建自己的注册表
	registry := prometheus.NewRegistry()
	// 使用我们自己的采集器
	registry.MustRegister(collector.NewNodeCollector())


	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry}))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error occur when start server %v", err)
	}
}
