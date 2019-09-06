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
	server := http.Server{
		Addr: ":10200",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			fmt.Fprintf(w, "[%s] go test server v2", time.Now().Format(time.RFC3339))
		}),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln("server start failed")
		}
	}()

	stopSig := make(chan os.Signal, 1)
	signal.Notify(stopSig, syscall.SIGINT, syscall.SIGTERM)
	<-stopSig

	log.Println("Shutdown-ing")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	graceful := make(chan struct{})
	go func() {
		defer close(graceful)
		if err := server.Shutdown(ctx); err != nil {
			log.Println("Shutdown error", err)
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("Shutdown timeout")
	case <-graceful:
		log.Println("bye")
	}
}
