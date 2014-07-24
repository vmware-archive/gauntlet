package main

import (
	"log"
	"net/http"
	"testing"

	cv "github.com/smartystreets/goconvey/convey"
)

// profiler: visit http://localhost:6060/debug/pprof for runtime info
func init() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}

func TestWebServerReturnsOnlyWhenItIsReachable(t *testing.T) {
	addr := "localhost:3000"
	cv.Convey("Given a call to StartWebServer", t, func() {
		cv.Convey("the webserver should be already listening on the chosen port when the function returns", func() {
			webserv := NewWebServer(addr)
			webserv.Start()
			defer webserv.Stop()
			cv.So(PortIsBound(addr), cv.ShouldEqual, true)
			//			url := "http://" + addr + "/results"
			//			resp, err := http.Get(url)
			//			fmt.Printf("http GET from '%s' returned err='%s' response='%s'\n", url, err, resp)
			//			cv.So(err, cv.ShouldEqual, nil)
		})
	})
}

func TestWebServerShutsdownWhenRequested(t *testing.T) {
	cv.Convey("Given a call to StartWebServer", t, func() {
		cv.Convey("the webserver be up after returning, and should terminate when requested", func() {
			addr := "localhost:3000"
			webserv := NewWebServer(addr)
			webserv.Start()
			cv.So(PortIsBound(addr), cv.ShouldEqual, true)
			webserv.Stop()
			cv.So(PortIsBound(addr), cv.ShouldEqual, false)

			// and again right away
			webserv = NewWebServer(addr)
			webserv.Start()
			cv.So(PortIsBound(addr), cv.ShouldEqual, true)
			webserv.Stop()
			cv.So(PortIsBound(addr), cv.ShouldEqual, false)

		})
	})
}
