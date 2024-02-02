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

func hello(w http.ResponseWriter, r *http.Request) {
	log.Printf("request %s from %s\n", r.RequestURI, r.Host)
	w.Write([]byte("ola world"))
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func main() {
	server := &http.Server{Addr: ":8080"}

	http.HandleFunc("/", hello)
	http.HandleFunc("/healthz", healthz)

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
