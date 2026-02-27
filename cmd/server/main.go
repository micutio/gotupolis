package main

// TODO: Make port and host configurable.
var (
	defaultHost = "localhost"
	defaultPort = 9001
)

func main() {
	tupleServer := NewServer(TCP, defaultHost, defaultPort)
	tupleServer.Launch()
}
