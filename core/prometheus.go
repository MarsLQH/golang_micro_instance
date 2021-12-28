package core

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	var (
		rpcDurations = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name:       "rpc_durations_seconds",
				Help:       "RPC latency  distributions",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			},
			[]string{"service"},
		)
	)

	//prometheus.NewHistogram(prometheus.HistogramOpts{Name:"rpc_duration_seconds",Help: "RPC latency",Buckets: prometheus.LinearBuckets(*normMean-5**normDomain, .5**normDomain, 20)})

	fmt.Println(rpcDurations)

}
