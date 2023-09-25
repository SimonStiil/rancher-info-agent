package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"path/filepath"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/client-go/util/homedir"
)

var (
	debug              bool
	prometheusEnabled  bool
	prometheusEndpoint string
	healthEndpoint     string
	kubeconfig         string
	port               string
	onlyRootEndpoint   bool
	requests           = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_endpoint_equests_count",
		Help: "The amount of requests to an endpoint",
	}, []string{"endpoint", "method"},
	)
	client KubeClient
)

type Health struct {
	Status string `json:"status"`
}

func HealthActuator(w http.ResponseWriter, r *http.Request) {
	if prometheusEnabled {
		requests.WithLabelValues(r.URL.EscapedPath(), r.Method).Inc()
	}
	if !(r.URL.Path == healthEndpoint) {
		log.Printf("@I %v %v %v - HealthActuator\n", r.Method, r.URL.Path, 404)
		http.NotFoundHandler().ServeHTTP(w, r)
		return
	}
	reply := Health{Status: "UP"}
	log.Printf("@I %v %v %v - HealthActuator\n", r.Method, r.URL.Path, 200)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reply)
	return
}

func MainHandler(w http.ResponseWriter, r *http.Request) {
	if prometheusEnabled {
		requests.WithLabelValues(r.URL.EscapedPath(), r.Method).Inc()
	}
	if onlyRootEndpoint {
		if !(r.URL.Path == "/") {
			log.Printf("@I %v %v %v - Main Handler\n", r.Method, r.URL.Path, 404)
			http.NotFoundHandler().ServeHTTP(w, r)
			return
		}
	}
	clusterList, err := client.CetClusters()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 Internal Server Error"))
		log.Printf("@I %v %v %v - Main Handler Request Error - %+v\n", r.Method, r.URL.Path, 500, err.Error())
		return
	}
	log.Printf("@I %v %v %v - Main Handler\n", r.Method, r.URL.Path, 200)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clusterList)
	return
}

func main() {
	flag.BoolVar(&debug, "debug", false, "Enable/(Disable) Debugging output")
	flag.StringVar(&port, "port", "8080", "port to use for the service (8080)")
	flag.BoolVar(&onlyRootEndpoint, "onlyRootEndpoint", false, "Enable/(Disable) Only reply json on / endpoint")
	flag.BoolVar(&prometheusEnabled, "prometheus", true, "(Enable)/Disable Prometheus endpoint")
	flag.StringVar(&prometheusEndpoint, "prometheusEndpoint", "/metrics", "custom prometheus endpoint (/metrics)")
	flag.StringVar(&healthEndpoint, "healthEndpoint", "/health", "custom health endpoint (/health)")

	if home := homedir.HomeDir(); home != "" {
		flag.StringVar(&kubeconfig, "kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	if debug {
		log.Println("@D Debugging enabled")
	}
	if prometheusEnabled {
		log.Printf("@I Metrics enabled at %v\n", prometheusEndpoint)
		http.Handle(prometheusEndpoint, promhttp.Handler())
	}

	client = KubeClient{Kubeconfig: kubeconfig, Debug: debug}
	http.HandleFunc(healthEndpoint, HealthActuator)
	http.HandleFunc("/", MainHandler)

	log.Printf("@I Serving on port %v\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
