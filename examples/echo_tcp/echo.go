package main

import (
	"dnet"
	"flag"
	"fmt"
	"log"
)

type echoServer struct {
	*dnet.EventServer
}

func (es *echoServer) OnInitComplete(srv dnet.Server) (action dnet.Action) {
	log.Printf("Echo server is listening on %s (multi-cores: %t, loops: %d)\n",
		srv.Addr.String(), srv.Multicore, srv.NumEventLoop)
	return
}

func (es *echoServer) React(frame []byte, c dnet.Conn) (out []byte, action dnet.Action) {
	// Echo synchronously.
	out = frame
	return
}

func main() {
	var port int
	var multicore bool

	// Example command: go run echo.go --port 9000 --multicore=true
	flag.IntVar(&port, "port", 9000, "--port 9000")
	flag.BoolVar(&multicore, "multicore", false, "--multicore true")
	flag.Parse()
	echo := new(echoServer)
	log.Fatal(dnet.Serve(echo, fmt.Sprintf("tcp://:%d", port), dnet.WithMulticore(multicore)))
}
