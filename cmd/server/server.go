package gotupolis

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

type Protocol int

const (
	TCP Protocol = 1
	UDP Protocol = 2
)

func (p Protocol) toString() string {
	switch p {
	case TCP:
		return "tcp"
	case UDP:
		return "udp"
	}
	return fmt.Sprintf("%d", int(p))
}

const (
	HOST = "localhost"
	PORT = "9001"
)

type Server struct {
	protocol Protocol
	host     string
	port     int
}

func (Server) init(protocol Protocol, host string, port int) Server {
	s := Server{protocol, host, port}
	return s
}

func handleIncomingRequest(conn net.Conn) {
	// store incoming data
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}
	// respond
	time := time.Now().Format("Monday, 02-Jan-06 15:04:05 MST")
	conn.Write([]byte("Hi back!\n"))
	conn.Write([]byte(time))

	// close conn
	conn.Close()
}

func (s Server) Launch() {
	fmt.Println("gotupolis server")
	listen, err := net.Listen(s.protocol.toString(), s.host+":"+fmt.Sprint(s.port))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	// close listener
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		go handleIncomingRequest(conn)
	}
}
