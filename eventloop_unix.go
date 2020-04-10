package dnet

import (
	"golang.org/x/sys/unix"
	"sync/atomic"
)

type eventloop struct {
	idx          int             // loop index in the server loops list
	svr          *server         // server in loop
	codec        ICodec          // codec for TCP
	packet       []byte          // read packet buffer
	poller       *netpoll.Poller // epoll or kqueue
	connCount    int32           // number of active connections in event-loop
	connections  map[int]*conn   // loop connections fd -> conn
	eventHandler EventHandler    // user eventHandler
}

func (el *eventloop) plusConnCount() {
	atomic.AddInt32(&el.connCount, 1)
}

func (el *eventloop) minusConnCount() {
	atomic.AddInt32(&el.connCount, -1)
}

func (el *eventloop) loadConnCount() int32 {
	return atomic.LoadInt32(&el.connCount)
}

func (el *eventloop) loopRun() {
	defer func() {
		if el.idx == 0 && el.svr.opts.Ticker {
			close(el.svr.ticktock)
		}
		el.svr.signalShutdown()
	}()

	if el.idx == 0 && el.svr.opts.Ticker {
		go el.loopTicker()
	}

	el.svr.logger.Printf("event-loop: %d exits with error: %v\n", el.idx, el.poller.Polling(el.handleEvent))
}

func (el *eventloop) loopAccept(fd int) error {
	if fd == el.svr.ln.fd {
		if el.svr.ln.pconn != nil {
			return el.loopReadUDP(fd)
		}
		nfd, sa, err := unix.Accept(fd)
		if err != nil {
			if err == unix.EAGAIN {
				return nil
			}
			return err
		}
		if err = unix.SetNonblock(nfd, true); err != nil {
			return err
		}
		c := newTCPConn(nfd, el, sa)
		if err = el.poller.AddRead(c.fd); err == nil {
			el.connections[c.fd] = c
			el.plusConnCount()
			return el.loopOpen(c)
		}
		return err
	}
	return nil
}

func (el *eventloop) loopOpen(c *conn) error {
	c.opened = true
	c.localAddr = el.svr.ln.lnaddr
}
