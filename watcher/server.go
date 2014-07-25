package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/glycerine/go-tigertonic"
)

type WebServer struct {
	Addr        string
	ServerReady chan bool // closed once server is listening on Addr
	Done        chan bool // recv on this to know that server is indeed shutdown
	LastReqBody string
	TigerSrv    *tigertonic.Server
}

func NewWebServer(addr string) *WebServer {

	s := &WebServer{
		Addr:        addr,
		ServerReady: make(chan bool),
		Done:        make(chan bool),
	}
	s.TigerSrv = tigertonic.NewServer(addr, s)

	return s
}

func (s *WebServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	buf := bytes.NewBuffer(nil)
	io.Copy(buf, r.Body)
	bodyAsString := string(buf.Bytes())
	fmt.Fprintf(w, "server got request body: '%s'\n", bodyAsString)
	fmt.Printf("server has bodyAsString: = %s\n", bodyAsString)

	s.LastReqBody = bodyAsString
}

func (s *WebServer) Start() *WebServer {
	fmt.Printf("\n\n WebServer::Start()\n")

	lsn, err := net.Listen("tcp", s.Addr)
	if nil != err {
		panic(err)
	}

	go s.TigerSrv.Serve(lsn)
	WaitUntilServerUp(s.Addr)
	close(s.ServerReady)
	return s
}

func (s *WebServer) Stop() {
	s.TigerSrv.Close()
	WaitUntilServerDown(s.Addr)
	close(s.Done)
}
