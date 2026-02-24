package server

import (
	"fmt"
	"log/slog"
	"net"
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

type Server struct {
	protocol Protocol
	host     string
	port     int
	logger   *slog.Logger
}

func NewServer(protocol Protocol, host string, port int) Server {
	s := Server{protocol, host, port, slog.Default()}
	return s
}

// TODO: handle actual tuple space client connections.
func (s *Server) handleIncomingRequest(conn net.Conn) {
	// store incoming data
	bufSize := 1024
	buffer := make([]byte, bufSize)
	_, readErr := conn.Read(buffer)
	if readErr != nil {
		s.logger.Error(
			"handleIncomingRequest: error reading from connection",
			slog.Any("error", readErr))
		return
	}
	// respond
	timeStamp := time.Now().Format("Monday, 02-Jan-06 15:04:05 MST")
	_, writeErr1 := conn.Write([]byte("Hi back!\n"))
	if writeErr1 != nil {
		s.logger.Error("handleIncomingRequest: error writing to connection",
			slog.Any("error", writeErr1))
		return
	}
	_, writeErr2 := conn.Write([]byte(timeStamp + "\n"))
	if writeErr2 != nil {
		s.logger.Error("handleIncomingRequest: error writing to connection",
			slog.Any("error", writeErr2))
		return
	}

	// close conn
	defer func() {
		closeErr := conn.Close()
		if closeErr != nil {
			s.logger.Error("handleIncomingRequest: error closing connection",
				slog.Any("error", closeErr))
		}
	}()
}

func (s *Server) Launch() {
	fmt.Println("gotupolis server")
	listen, listenErr := net.Listen(s.protocol.toString(), s.host+":"+fmt.Sprint(s.port))
	if listenErr != nil {
		s.logger.Error("Launch: error listen to host:port",
			slog.String("host", s.host),
			slog.Int("port", s.port),
			slog.Any("error", listenErr))
	}
	// close listener
	defer func() {
		closeErr := listen.Close()
		if closeErr != nil {
			s.logger.Error("Launch: error closing Listener", slog.Any("error", closeErr))
		}
	}()

	for {
		conn, connErr := listen.Accept()
		if connErr != nil {
			s.logger.Error("Launch: error accepting connection", slog.Any("error", connErr))
		}
		go s.handleIncomingRequest(conn)
	}
}
