package main

// Import necessary packages
import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	oneMb = 1024 * 1024
)

var helloCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "hello_request_count",
		Help: "No of request handled by Hello handler",
	},
)

var shouldFailOnHealthz = false

func hello(w http.ResponseWriter, r *http.Request) {
	helloCounter.Inc()
	log.Printf("request %s from %s\n", r.RequestURI, r.Host)
	w.Write([]byte("hello world"))
}

func notfound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	if shouldFailOnHealthz {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}

func heavy(w http.ResponseWriter, r *http.Request) {
	log.Printf("request %s from %s\n", r.RequestURI, r.Host)
	v := make([]int, 10*oneMb)
	v[0] = 0x00
	w.Write([]byte("heavy request"))
}

func init() {
	_, shouldFailOnHealthz = os.LookupEnv("FLAKY")
	if shouldFailOnHealthz {
		log.Println("whops, i might be flaky O.o")
	}
}

func main() {
	prometheus.MustRegister(helloCounter)

	server := &http.Server{Addr: ":8080"}
	http.HandleFunc("/", notfound)
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/heavy", heavy)
	http.HandleFunc("/healthz", healthz)
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/flaky", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("i just got flaky :-/"))
		shouldFailOnHealthz = true
	})

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("Booting up server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %s", err)
		}
	}()

	<-stopChan
	log.Println("Shutting down server")
	server.Shutdown(context.Background())
}
