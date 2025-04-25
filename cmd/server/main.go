package main

import "github.com/micutio/gotupolis/pkg/server"

func main() {
	tupleServer := server.NewServer(server.TCP, "localhost", 9001)
	tupleServer.Launch()
}
