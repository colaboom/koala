package middleware

import (
	"github.com/koala/middleware/prometheus"
)

var (
	DefaultServerMetrics = prometheus.NewServerMetrics()
)