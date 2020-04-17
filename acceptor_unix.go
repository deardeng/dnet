package dnet

import (
	"golang.org/x/sys/unix"
	"log"
)

func (svr *server) acceptNewConnection(fd int) error {
	nfd, sa, err := unix.Accept(fd)
	log.Println("accept...")
	if err != nil {
		if err != unix.EAGAIN {
			return nil
		}
		return err
	}
	if err := unix.SetNonblock(nfd, true); err != nil {
		return err
	}
	el := svr.subLoopGroup.next(nfd)
	c := newTCPConn(nfd, el, sa)
	_ = el.poller.Trigger(func() (err error) {
		if err = el.poller.AddRead(nfd); err != nil {
			return
		}
		el.connections[nfd] = c
		el.plusConnCount()
		err = el.loopOpen(c)
		return
	})
	return nil
}
