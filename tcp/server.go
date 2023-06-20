package tcp

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/hattorious/echoserver/version"
)

// <remote_IP_address> - [<timestamp>] "- - <request_protocol>" - <request_bytes>
const tcpLogLine = "%s - - [%s] \"- - TCP\" - %v"

func StartTCPServer(port int, verbose bool) {
	addr := &net.TCPAddr{
		Port: port,
	}

	server, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("INFO: listening on %s; TCP", addr)
	defer server.Close()

	for {
		conn, err := server.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}

		go handleTcpConnection(conn, verbose)
	}
}

func handleTcpConnection(conn *net.TCPConn, verbose bool) {
	defer conn.Close()

	b := make([]byte, 1024)
	for {
		n, err := conn.Read(b)
		if err != nil {
			if err == io.EOF {
				// client closed conn
				return
			}

			log.Printf("DEBUG: net.TCPConn.Read() error: %s", err)
			break
		}

		resp := fmt.Sprintf("%s %v > %s", version.ServerID, n, b[:n])

		_, err = conn.Write([]byte(resp))
		if err != nil {
			if verbose {
				log.Printf("DEBUG: net.TCPConn.Write() error: %s", err)
			}
		}

		if verbose {
			log.Printf(tcpLogLine, conn.RemoteAddr(), time.Now().Format("02/Jan/2006:15:04:05 -0700"), n)
		}
	}
}
