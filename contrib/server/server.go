package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	webroot = flag.String("webroot", ".", "Directory to serve files from")
	port    = flag.String("port", "8090", "Port to serve on")
)

func main() {
	flag.Parse()

	fs := http.FileServer(http.Dir(*webroot))
	server := &http.Server{
		Addr: ":" + *port,
		Handler: http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			resp.Header().Add("Cache-Control", "no-cache")
			if strings.HasSuffix(req.URL.Path, ".wasm") {
				resp.Header().Set("content-type", "application/wasm")
			}
			fs.ServeHTTP(resp, req)
		}),
	}

	// Handle graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Serving %s on http://localhost:%s", *webroot, *port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error: %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down server...")
}
