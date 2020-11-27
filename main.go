package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Println(req.Method, req.URL.Path, req.RemoteAddr)
		_, _ = io.Copy(w, req.Body)
	})

	handler.HandleFunc("/version", func(w http.ResponseWriter, _ *http.Request) {
		const version = "latest"
		_, _ = fmt.Fprintln(w, version)
	})

	handler.HandleFunc("/ping", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprintln(w, "pong")
	})

	server := http.Server{Addr: ":8080", Handler: handler}
	go func() {
		stopSig := make(chan os.Signal, 1)
		signal.Notify(stopSig, syscall.SIGINT, syscall.SIGTERM)
		<-stopSig
		log.Println("graceful stopping")
		_ = server.Shutdown(context.TODO())
	}()

	log.Println("listen and serving")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalln("server start failed", err)
	}
}
