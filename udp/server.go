package udp

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/hattorious/echoserver/version"
)

// <remote_IP_address> - [<timestamp>] "- - <request_protocol>" - <request_bytes>
const udpLogLine = "%s - - [%s] \"- - UDP\" - %v"

func StartUDPServer(port int, verbose bool) {
	addr := &net.UDPAddr{
		Port: port,
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("INFO: listening on %s; UDP", addr)
	defer conn.Close()

	b := make([]byte, 2048)

	for {
		n, remote, err := conn.ReadFromUDP(b)
		if err != nil {
			if verbose {
				log.Printf("DEBUG: net.ReadFromUDP() error: %s", err)
			}
			continue
		}

		go handleUdpPacketConn(conn, remote, b[:n], n, verbose)
	}
}

func handleUdpPacketConn(conn net.PacketConn, addr net.Addr, buf []byte, n int, verbose bool) {
	resp := fmt.Sprintf("%s %v > %s", version.ServerID, n, buf)

	_, err := conn.WriteTo([]byte(resp), addr)
	if err != nil {
		if verbose {
			log.Printf("DEBUG: net.PacketConn.WriteTo() error: %s", err)
		}
	}

	if verbose {
		log.Printf(udpLogLine, addr, time.Now().Format("02/Jan/2006:15:04:05 -0700"), n)
	}
}
