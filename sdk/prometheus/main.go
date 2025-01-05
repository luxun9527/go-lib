package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	dto "github.com/prometheus/client_model/go"
)

// 定义 Prometheus 指标
var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})

	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"method", "endpoint", "status"})

	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests",
		Buckets: prometheus.DefBuckets, // 默认的桶分布
	}, []string{"method", "endpoint", "status"})

	httpRequestsRate = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "http_requests_rate",
		Help: "Average number of HTTP requests per second",
	}, []string{"endpoint"})
)

// 中间件：统计请求数和延迟
func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 包装 ResponseWriter 以捕获状态码
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		// 记录请求总数
		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(rw.status)).Inc()

		// 记录请求延迟
		duration := time.Since(start).Seconds()
		httpRequestDuration.WithLabelValues(r.Method, r.URL.Path, http.StatusText(rw.status)).Observe(duration)
	})
}

// 自定义 ResponseWriter 以捕获状态码
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// 模拟处理逻辑
func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

// 计算每秒请求速率
func calculateRequestRate() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// 获取所有接口的请求总数
		metrics, err := httpRequestsTotal.GetMetricWithLabelValues("GET", "/", "OK")
		if err != nil {
			log.Printf("Failed to get metrics: %v", err)
			continue
		}
		// 获取当前计数器的值
		counter := &dto.Metric{}
		if err := metrics.Write(counter); err != nil {
			log.Printf("Failed to write metrics: %v", err)
			continue
		}

		// 计算每秒请求速率
		rate := counter.GetCounter().GetValue()
		httpRequestsRate.WithLabelValues("/").Set(rate)
	}
}

func main() {
	// 启动模拟指标
	recordMetrics()

	// 启动计算每秒请求速率的 goroutine
	go calculateRequestRate()

	// 创建一个简单的 HTTP 处理函数
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))

	})

	// 使用中间件包装处理函数
	wrappedHandler := metricsMiddleware(handler)

	// 暴露 Prometheus 指标
	http.Handle("/metrics", promhttp.Handler())

	// 启动 HTTP 服务器
	http.Handle("/", wrappedHandler)
	http.Handle("/test", wrappedHandler)
	http.Handle("/test1", wrappedHandler)
	log.Println("Server started on :9999")
	http.ListenAndServe(":9999", nil)
}
