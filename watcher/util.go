package main

import (
	"fmt"
	"net"
	"time"
)

func WaitUntilServerUp(addr string) {
	attempt := 1
	for {
		if PortIsBound(addr) {
			return
		}

		time.Sleep(50 * time.Millisecond)
		attempt++
		if attempt > 40 {
			panic(fmt.Sprintf("could not connect to server at '%s' after 40 tries of 50msec", addr))
		}
	}
}

func WaitUntilServerDown(addr string) {
	attempt := 1
	for {
		if !PortIsBound(addr) {
			return
		}
		time.Sleep(50 * time.Millisecond)
		attempt++
		if attempt > 40 {
			panic(fmt.Sprintf("could always connect to server at '%s' after 40 tries of 50msec", addr))
		}
	}
}

func PortIsBound(addr string) bool {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
