package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var fooCount = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "foo_total",
	Help: "Number of foo successfully processed.",
})

var hits = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "hits",
}, []string{"status", "path"})

func main() {
	prometheus.MustRegister(fooCount, hits)

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		hits.WithLabelValues("200", r.URL.String()).Inc()
		fooCount.Add(1)
		fmt.Fprintf(w, "foo_total increased")
	})

	http.HandleFunc("/500", func(w http.ResponseWriter, r *http.Request) {
		hits.WithLabelValues("500", r.URL.String()).Inc()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
