package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

// The golang net/http standard library seems to insist
// that we only run one webserver using ListenAndServer(),
// as there is no way to shutdown the first or start a
// second web server. Yuck. But for now, we just go with it.
//
var SingletonWebServer *WebServer
var SingletonWebServerStarted bool

type WebServer struct {
	Addr        string
	ServerReady chan bool // closed once server is listening on Addr
	RequestStop chan bool // send on this to tell server to shutdown
	Done        chan bool // recv on this to know that server is indeed shutdown
	LastReqBody string
}

//func (s *WebServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	fmt.Printf("WebServer::ServeHTTP running ...\n")
//	s.ReceivedReq <- r
//}

func NewWebServer(addr string) *WebServer {

	if SingletonWebServer != nil {
		return SingletonWebServer
	}

	SingletonWebServer = &WebServer{
		Addr:        addr,
		ServerReady: make(chan bool),
		RequestStop: make(chan bool),
		Done:        make(chan bool),
	}

	return SingletonWebServer
}

func (webserv *WebServer) Start() {

	if SingletonWebServerStarted {
		return
	}
	SingletonWebServerStarted = true

	//fmt.Printf("\n top of StartWebServer\n")

	resultsHandler := func(w http.ResponseWriter, r *http.Request) {
		//fmt.Printf("resultsHandler running ...\n")

		buf := bytes.NewBuffer(nil)
		io.Copy(buf, r.Body)
		bodyAsString := string(buf.Bytes())
		fmt.Fprintf(w, "server got request body: '%s'\n", bodyAsString)
		//fmt.Printf("server has bodyAsString: = %s\n", bodyAsString)

		webserv.LastReqBody = bodyAsString
	}

	go func() {
		http.HandleFunc("/results", resultsHandler)
		fmt.Printf("listening on %s and responding to /results\n", webserv.Addr)
		log.Fatal(http.ListenAndServe(webserv.Addr, nil))
	}()

	WaitUntilServerUp(webserv.Addr)
	close(webserv.ServerReady)
}

func (s *WebServer) Stop() {

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

func PortIsBound(addr string) bool {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
