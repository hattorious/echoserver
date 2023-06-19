package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/hattorious/echoserver/http"
	"github.com/hattorious/echoserver/udp"
)

var (
	verbose bool
	ports   struct {
		http1 string
		http2 string
		tcp   int
		udp   int
	}
)

func init() {
	flag.BoolVar(&verbose, "verbose", false, "enable verbose logging")
	flag.StringVar(&ports.http1, "http1", getStringEnv("ECHO_HTTP1_PORT", "8001"), "give me a port number for HTTP/1")
	flag.StringVar(&ports.http2, "http2", getStringEnv("ECHO_HTTP2_CLEARTEXT_PORT", "8002"), "give me a port number for HTTP/2 over clear-text")
	flag.IntVar(&ports.udp, "udp", getIntEnv("ECHO_UDP_PORT", 8007), "give me a port number for UDP")
}

func main() {
	flag.Parse()

	go http.StartHttpServer(ports.http1, verbose)
	go http.StartHttp2CleartextServer(ports.http2, verbose)
	udp.StartUDPServer(ports.udp, verbose)
}

func getStringEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getIntEnv(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	valueInt, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("ERROR: couldn't decode environment variable %s=%s to a port number; %s", key, value, err)
	}
	return valueInt
}
