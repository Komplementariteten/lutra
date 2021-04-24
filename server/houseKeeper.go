package server

import (
	"context"
	"fmt"
	"time"
)

type houseKeeperControl struct {
	e chan bool
	q chan bool
	s *HTTPServer
}
type HouseKeeperTimer struct {
	t    time.Duration
	h    *houseKeeperControl
	Quit chan bool
}

type Cleanabel interface {
	Cleanup(ctx context.Context)
	ShouldBeExecuted(n int) bool
}

func StartNewHouseKeeper(server *HTTPServer, trigger time.Duration) *HouseKeeperTimer {
	ctrl := &houseKeeperControl{
		e: make(chan bool),
		q: make(chan bool),
		s: server,
	}
	go func(s *houseKeeperControl) {
		count := 0
		ctx, cancel := context.WithCancel(context.Background())
		for {
			select {
			case <-s.e:
				fmt.Println("Execute housekeeping")
				if s.s.i.ShouldBeExecuted(count) {
					s.s.i.Cleanup(ctx)
				}
				if s.s.a.ShouldBeExecuted(count) {
					s.s.a.Cleanup(ctx)
				}
				if s.s.p.ShouldBeExecuted(count) {
					s.s.p.Cleanup(ctx)
				}
				count++
			case <-s.q:
				cancel()
				close(s.e)
				close(s.q)
				s.s.Shutdown()
				fmt.Println("Server shutdown")
				return
			}
		}
	}(ctrl)
	t := &HouseKeeperTimer{
		Quit: make(chan bool),
		t:    trigger,
		h:    ctrl,
	}
	go func(t *HouseKeeperTimer) {
		for {
			select {
			case <-t.Quit:
				t.h.e <- true
				t.h.q <- true
				fmt.Println("Request to quit")
				close(t.Quit)
				return
			case <-time.After(t.t):
				t.h.e <- true
			}
		}
	}(t)
	return t
}
