package main

// Import necessary packages
import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var shouldFailOnHealthz = false

func hello(w http.ResponseWriter, r *http.Request) {
	log.Printf("request %s from %s\n", r.RequestURI, r.Host)
	w.Write([]byte("ola world"))
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

func init() {
	_, shouldFailOnHealthz = os.LookupEnv("FLAKY")
	if shouldFailOnHealthz {
		log.Println("whops, i might be flaky O.o")
	}
}

func main() {
	server := &http.Server{Addr: ":8080"}
	http.HandleFunc("/", notfound)
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/healthz", healthz)
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
