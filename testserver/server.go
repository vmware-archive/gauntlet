package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)


func resultsHandler(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\n\n top of resultsHandler: request r = %#v\n", r)
		buf := bytes.NewBuffer(nil)
		io.Copy(buf, r.Body)

		fmt.Printf("server debug: request.Method = %s\n", r.Method)
		fmt.Printf("server debug: request r.Body = '%s'\n", buf)

		fmt.Fprintf(w, "server got request body:, %s", string(buf.Bytes()))

		switch r.Method {
		   case "POST": return createResult(w,r)
		   case "GET":  return listResult(w,r)
		}

}

func createResult(w http.ResponseWriter, r *http.Request) {

}
func listResult(w http.ResponseWriter, r *http.Request) {

}

func main() {

	http.HandleFunc("/results", resultsHandler)
	addr := "localhost:3000"
	fmt.Printf("listening on %s and responding to /results\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
