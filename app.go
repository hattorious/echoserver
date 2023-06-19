package main

import (
	"flag"
	"os"

	"github.com/hattorious/echoserver/http"
)

var (
	verbose bool
	ports   struct {
		http1 string
		http2 string
	}
)

func init() {
	flag.BoolVar(&verbose, "verbose", false, "enable verbose logging")
	flag.StringVar(&ports.http1, "http1-port", getEnv("ECHO_HTTP1_PORT", "8001"), "give me a port number for HTTP/1")
	flag.StringVar(&ports.http2, "http2-port", getEnv("ECHO_HTTP2_CLEARTEXT_PORT", "8002"), "give me a port number for HTTP/2 over clear-text")
}

func main() {
	flag.Parse()

	go http.StartHttpServer(ports.http1, verbose)
	http.StartHttp2CleartextServer(ports.http2, verbose)
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
