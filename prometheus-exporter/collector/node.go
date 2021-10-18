package collector

import (
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

var hostname string

type NodeCollector struct {
	requestDesc    *prometheus.Desc //Counter类型
	goroutinesDesc *prometheus.Desc //Gauge
	threadsDesc    *prometheus.Desc //Gauge
	summaryDesc    *prometheus.Desc //Summary
	histogramDesc  *prometheus.Desc //Histogram
	nodeMetrics    nodeStatsMetrics //混合方式
}

type nodeStatsMetrics []struct {
	desc    *prometheus.Desc
	eval    func(*mem.VirtualMemoryStat) float64
	valType prometheus.ValueType
}

func NewNodeCollector() prometheus.Collector {
	host, _ := host.Info()
	hostname = host.Hostname
	return &NodeCollector{
		requestDesc: prometheus.NewDesc(
			"total_request_total",
			"request总请求数",
			[]string{"DYNAMIC_HOST_NAME"},                                           // 变量标签，标签值可变，在metrics中设置值
			prometheus.Labels{"static_label1": "value1", "static_label2": "value2"}, // 静态标签
		),

		goroutinesDesc: prometheus.NewDesc(
			"goroutines__nums",
			"goroutine协程数",
			nil,
			nil,
		),

		threadsDesc: prometheus.NewDesc(
			"thread_nums",
			"thread数",
			nil,
			nil,
		),

		summaryDesc: prometheus.NewDesc(
			"summary_http_request_duration_seconds", // http请求持续时间
			"summary类型",
			[]string{"code", "method"},
			prometheus.Labels{"static_label1": "static_value1"},
		),

		histogramDesc: prometheus.NewDesc(
			"histogram_http_request_duration_seconds",
			"histogram类型",
			[]string{"code", "method"},
			prometheus.Labels{"static_label1": "static_value1"},
		),

		nodeMetrics: nodeStatsMetrics{
			{
				desc: prometheus.NewDesc(
					"total_mem",
					"内存总量",
					nil,
					nil,
				),
				eval:    func(ms *mem.VirtualMemoryStat) float64 { return float64(ms.Total) / 1e9 },
				valType: prometheus.GaugeValue,
			},
			{
				desc: prometheus.NewDesc(
					"free_mem",
					"内存空闲",
					nil,
					nil,
				),
				eval:    func(ms *mem.VirtualMemoryStat) float64 { return float64(ms.Free) / 1e9 },
				valType: prometheus.GaugeValue,
			},
		},
	}
}

//实现采集器Describe接口
func (n *NodeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- n.requestDesc
	ch <- n.goroutinesDesc
	ch <- n.threadsDesc
	for _, metrics := range n.nodeMetrics {
		ch <- metrics.desc
	}
}

//实现采集器Collect接口,真正采集动作
func (n *NodeCollector) Collect(ch chan<- prometheus.Metric) {
	nums := 0
	for i := 0; i < 10; i++ {
		nums++
	}
	ch <- prometheus.MustNewConstMetric(n.requestDesc, prometheus.CounterValue, float64(nums), hostname)

	ch <- prometheus.MustNewConstMetric(n.goroutinesDesc, prometheus.GaugeValue, float64(runtime.NumGoroutine()))

	v, _ := runtime.ThreadCreateProfile(nil)
	ch <- prometheus.MustNewConstMetric(n.threadsDesc, prometheus.GaugeValue, float64(v))

	ch <- prometheus.MustNewConstSummary(n.summaryDesc,
		4711,
		403.34,
		map[float64]float64{0.5: 42.3, 0.9: 32.3},
		"200", "get",
	)

	vm, _ := mem.VirtualMemory()
	for _, metrics := range n.nodeMetrics {
		ch <- prometheus.MustNewConstMetric(metrics.desc, metrics.valType, metrics.eval(vm))
	}
}
