package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "v1")
	})

	http.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "pong")
	})

	server := http.Server{Addr: ":10200"}
	go func() {
		log.Println("start the server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln("server start failed", err)
		}
	}()

	stopSig := make(chan os.Signal, 1)
	signal.Notify(stopSig, syscall.SIGINT, syscall.SIGTERM)
	<-stopSig

	log.Println("graceful stopping")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	graceful := make(chan struct{})
	go func() {
		defer close(graceful)
		if err := server.Shutdown(ctx); err != nil {
			log.Println("shutdown error", err)
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("shutdown timeout")
	case <-graceful:
		log.Println("bye")
	}
}
