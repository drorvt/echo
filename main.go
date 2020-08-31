package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	const version = "v1"

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Println(req.Method, req.URL.Path, req.RemoteAddr)
		_, _ = io.Copy(w, req.Body)
	})

	http.HandleFunc("/v", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(version))
	})

	server := http.Server{Addr: ":8080"}
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
