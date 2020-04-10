package dnet

import (
	"net"
	"os"
	"sync"
)

type listener struct {
	f             *os.File
	fd            int
	ln            net.Listener
	once          sync.Once
	pconn         net.PacketConn
	lnaddr        net.Addr
	addr, network string
}
