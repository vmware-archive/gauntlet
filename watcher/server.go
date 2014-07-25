package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
	"github.com/glycerine/go-tigertonic"

	_ "net/http/pprof"
)

func exampleTigerTonic() *tigertonic.Server {
	s := tigertonic.NewServer(
		"127.0.0.1:3000",
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}),
	)
	l1, err := net.Listen("tcp", s.Addr)
	if nil != err {
		panic(err)
	}
	go s.Serve(l1)
	if _, err := http.Get(fmt.Sprintf("http://%s", l1.Addr())); nil != err {
		panic(err)
	}
	return s
}

type WebServer struct {
	Addr        string
	ServerReady chan bool // closed once server is listening on Addr
	RequestStop chan bool // close this to tell server to shutdown
	Done        chan bool // recv on this to know that server is indeed shutdown
	LastReqBody string
	Lsn         net.Listener
	Tts         *tigertonic.Server
	Mux         *http.ServeMux
	Started     bool
}

func NewWebServer(addr string) *WebServer {

	ws := &WebServer{
		Addr:        addr,
		ServerReady: make(chan bool),
		RequestStop: make(chan bool),
		Done:        make(chan bool),
	}
	ws.Tts = tigertonic.NewServer(addr, ws)

	lsn, err := net.Listen("tcp", ws.Addr)
	if nil != err {
		panic(err)
	}
	ws.Lsn = lsn
	return ws
}

func (webserv *WebServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("resultsHandler (on %p) for %s running ...\n", webserv, webserv.Addr)

	//	if webserv.StopRequested() {
	//		panic("detected resultsHandler() call on already stop-requested web-server")
	//	}

	buf := bytes.NewBuffer(nil)
	io.Copy(buf, r.Body)
	bodyAsString := string(buf.Bytes())
	fmt.Fprintf(w, "server got request body: '%s'\n", bodyAsString)
	fmt.Printf("server %p has bodyAsString: = %s\n", webserv, bodyAsString)

	webserv.LastReqBody = bodyAsString
}

func (s *WebServer) Start() *WebServer {
	fmt.Printf("\n\n WebServer::Start() for %p\n", s)

	go s.Tts.Serve(s.Lsn)
	WaitUntilServerUp(s.Addr)
	close(s.ServerReady)
	return s
}

func (s *WebServer) Stop() {
	s.Tts.Close()
	WaitUntilServerDown(s.Addr)
	close(s.Done)
}

// determine if stop has been requested by
// checking if s.RequestStop has been closed.
func (s *WebServer) StopRequested() bool {
	select {
	case <-s.RequestStop:
		return true
	default:
		return false
	}
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
