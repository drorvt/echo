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
	file, err := os.OpenFile("/log/k8s.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Open log file error")
	}
	defer file.Close()

	log.SetOutput(file)
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Println(req.Method, req.URL.Path, req.RemoteAddr)
		_, _ = io.Copy(w, req.Body)
	})

	server := http.Server{Addr: ":10200"}
	go func() {
		log.Println("listen and serve...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln("server start failed", err)
		}
	}()

	stopSig := make(chan os.Signal, 1)
	signal.Notify(stopSig, syscall.SIGINT, syscall.SIGTERM)
	<-stopSig

	log.Println("graceful stopping...")
	server.Shutdown(context.TODO())
}
