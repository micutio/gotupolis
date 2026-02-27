package main

import "github.com/micutio/gotupolis/cmd/server"

// TODO: Make port and host configurable.
var (
	defaultHost = "localhost"
	defaultPort = 9001
)

func main() {
	tupleServer := server.NewServer(server.TCP, defaultHost, defaultPort)
	tupleServer.Launch()
}
