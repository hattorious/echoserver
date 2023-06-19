package http

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/hattorious/echoserver/version"
	"github.com/netinternet/remoteaddr"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func indentFPrintln(w io.Writer, level int, format string, args ...interface{}) {
	var b strings.Builder

	for i := 0; i < level; i++ {
		fmt.Fprintf(&b, "\t")
	}

	fmt.Fprintf(&b, format, args...)

	_, _ = fmt.Fprintln(w, b.String())
}

func handle(next http.HandlerFunc, verbose bool) http.Handler {
	if !verbose {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next(w, r)

		// <remote_IP_address> - [<timestamp>] "<request_method> <request_path> <request_protocol>" -
		log.Printf("%s - - [%s] \"%s %s %s\" - -", r.RemoteAddr, time.Now().Format("02/Jan/2006:15:04:05 -0700"), r.Method, r.URL.Path, r.Proto)
	})
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Server", version.ServerID)

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	remoteIP, _ := remoteaddr.Parse().AddHeaders([]string{"X-Envoy-External-Address"}).IP(req)

	reqDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Fatal(err)
	}

	envvars := os.Environ()
	sort.Strings(envvars)

	indentFPrintln(w, 0, "%s", hostname)
	indentFPrintln(w, 0, "v%s", version.String())
	indentFPrintln(w, 0, "")
	indentFPrintln(w, 0, "REQUEST:")
	indentFPrintln(w, 0, "--------")
	indentFPrintln(w, 0, "remote address: %s", req.RemoteAddr)
	indentFPrintln(w, 0, "remote IP: %s", remoteIP)
	indentFPrintln(w, 0, "URI: %s", req.RequestURI)
	indentFPrintln(w, 0, "version: %s", req.Proto)
	indentFPrintln(w, 0, "TLS: %v", req.TLS != nil)
	indentFPrintln(w, 0, "content length: %v", req.ContentLength)
	indentFPrintln(w, 0, "")
	indentFPrintln(w, 0, "%s", string(reqDump))
	indentFPrintln(w, 0, "ENVIRONMENT:")
	indentFPrintln(w, 0, "------------")

	for _, env := range envvars {
		indentFPrintln(w, 0, "%s", env)
	}
}

func StartHttpServer(port string, verbose bool) {
	server := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: handle(indexHandler, verbose),
	}

	log.Printf("listening on %s; HTTP/1, HTTP/1.1", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func StartHttp2CleartextServer(port string, verbose bool) {
	server := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: h2c.NewHandler(handle(indexHandler, verbose), &http2.Server{}),
	}

	log.Printf("listening on %s; HTTP/1, HTTP/1.1, h2c (clear-text HTTP/2, upgrade, prior knowledge)", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func StartHttp2TLSServer(port string, verbose bool) {
	server := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: handle(indexHandler, verbose),
	}

	err := http2.ConfigureServer(server, &http2.Server{})

	if err != nil {
		panic(err)
	}

	log.Printf("listening on %s; HTTP/1, HTTP/2", server.Addr)
	if err := server.ListenAndServeTLS("./default.pem", "./default.key"); err != nil {
		panic(err)
	}
}
