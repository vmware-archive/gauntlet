package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	_ "net/http/pprof"
)

// profiler: visit http://localhost:6060/debug/pprof for runtime info
/*
func init() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}
*/

/*
Portions of the code in this file are derived from source code that is licensed as follows:

writeup of stoppableListener approach:
http://www.hydrogen18.com/blog/stop-listening-http-server-go.html


Copyright (c) 2014, Eric Urban
All rights reserved.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

type WebServer struct {
	Addr        string
	ServerReady chan bool // closed once server is listening on Addr
	RequestStop chan bool // send on this to tell server to shutdown
	Done        chan bool // recv on this to know that server is indeed shutdown
	LastReqBody string
	OL          net.Listener
	SL          *StoppableListener
	Server      http.Server
	Mux         *http.ServeMux
}

func NewWebServer(addr string) *WebServer {

	originalListener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	stoppableListener, err := NewSL(originalListener)
	if err != nil {
		panic(err)
	}

	ws := &WebServer{
		Addr:        addr,
		ServerReady: make(chan bool),
		RequestStop: make(chan bool),
		Done:        make(chan bool),
		OL:          originalListener,
		SL:          stoppableListener,
	}
	mux := http.NewServeMux()
	ws.Server.Handler = mux
	ws.Mux = mux
	return ws
}

func (webserv *WebServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("resultsHandler (on %p) for %s running ...\n", webserv, webserv.Addr)

	buf := bytes.NewBuffer(nil)
	io.Copy(buf, r.Body)
	bodyAsString := string(buf.Bytes())
	fmt.Fprintf(w, "server got request body: '%s'\n", bodyAsString)
	fmt.Printf("server %p has bodyAsString: = %s\n", webserv, bodyAsString)

	webserv.LastReqBody = bodyAsString
}

func (webserv *WebServer) Start() *WebServer {

	//fmt.Printf("\n top of StartWebServer\n")
	/*
		resultsHandler := func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("resultsHandler for %s running ...\n", webserv.Addr)

			buf := bytes.NewBuffer(nil)
			io.Copy(buf, r.Body)
			bodyAsString := string(buf.Bytes())
			fmt.Fprintf(w, "server got request body: '%s'\n", bodyAsString)
			fmt.Printf("server %p has bodyAsString: = %s\n", webserv, bodyAsString)

			webserv.LastReqBody = bodyAsString
		}
	*/
	//webserv.Mux.HandleFunc("/results", resultsHandler)
	webserv.Mux.Handle("/results", webserv)

	go func(ws *WebServer) {

		fmt.Printf("listening on %s and responding to /results\n", ws.Addr)
		// blocks until webserv.SL.Stop()
		ws.Server.Serve(webserv.SL)
		fmt.Printf("\n\n 88888888888 !!! webserv( %p  ).Server.Serve() has returned !!! 888888888888 \n\n", ws)

		close(webserv.Done)
	}(webserv)

	WaitUntilServerUp(webserv.Addr)
	close(webserv.ServerReady)
	return webserv
}

func (s *WebServer) Stop() {
	s.SL.Stop()
	fmt.Printf("\n\n webserv::Stop() request sent, on %p \n\n", s)
	<-s.Done
	WaitUntilServerDown(s.Addr)
}

func WaitUntilServerUp(addr string) {
	attempt := 1
	for {
		if PortIsBound(addr) {
			return
		}
		//fmt.Printf("WaitUntilServerUp: on attempt %d, sleep then try again\n", attempt)
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
		//fmt.Printf("WaitUntilServerUp: on attempt %d, sleep then try again\n", attempt)
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

type StoppableListener struct {
	*net.TCPListener          //Wrapped listener
	stop             chan int //Channel used only to indicate listener should shutdown
}

func NewSL(l net.Listener) (*StoppableListener, error) {
	tcpL, ok := l.(*net.TCPListener)

	if !ok {
		return nil, errors.New("Cannot wrap listener")
	}

	retval := &StoppableListener{}
	retval.TCPListener = tcpL
	retval.stop = make(chan int)

	return retval, nil
}

var StoppedError = errors.New("Listener stopped")

func (sl *StoppableListener) Accept() (net.Conn, error) {

	for {
		//Wait up to 100msec for a new connection
		sl.SetDeadline(time.Now().Add(100 * time.Millisecond))

		newConn, err := sl.TCPListener.Accept()

		//Check for the channel being closed
		select {
		case <-sl.stop:
			return nil, StoppedError
		default:
			//If the channel is still open, continue as normal
		}

		if err != nil {
			netErr, ok := err.(net.Error)

			//If this is a timeout, then continue to wait for
			//new connections
			if ok && netErr.Timeout() && netErr.Temporary() {
				continue
			}
		}

		return newConn, err
	}
}

func (sl *StoppableListener) Stop() {
	close(sl.stop)
}
