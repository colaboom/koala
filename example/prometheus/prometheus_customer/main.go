package main

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"
)

var (
	addr              = flag.String("listen-address", ":8080", "The address to listen on for HTTP request.")
	uniformDomain     = flag.Float64("uniform.domain", 0.0002, "The domain for the uniform distribution.")
	normDomain        = flag.Float64("normal.domain", 0.0002, "The domain for the normal distribution.")
	normMean          = flag.Float64("normal.mean", 0.00001, "The mean for the normal distribution.")
	oscillationPeriod = flag.Duration("oscillation-period", 10*time.Minute, "The duration of the rate")
)

var (
	rpcDurations = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "rpc_duraions_seconds",
			Help:       "RPC latency distributions",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"service"},
	)
)

func init() {
	prometheus.MustRegister(rpcDurations)
	prometheus.MustRegister(collectors.NewBuildInfoCollector())
}

func main() {
	flag.Parse()
	start := time.Now()
	oscillationPeriod := func() float64 {
		return 2 + math.Sin(2*math.Pi*float64(time.Since(start))/float64(*oscillationPeriod))
	}

	go func() {
		for {
			v := rand.Float64() * *uniformDomain
			rpcDurations.WithLabelValues("uniform").Observe(v)
			time.Sleep(time.Duration(100*oscillationPeriod()) * time.Millisecond)
		}
	}()

	go func() {
		for {
			v := (rand.Float64() * *normDomain) + *normMean
			rpcDurations.WithLabelValues("normal").Observe(v)
			time.Sleep(time.Duration(75*oscillationPeriod()) * time.Millisecond)
		}
	}()

	go func() {
		for {
			v := rand.ExpFloat64() / 1e6
			rpcDurations.WithLabelValues("exponential").Observe(v)
			time.Sleep(time.Duration(50*oscillationPeriod()) * time.Millisecond)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
