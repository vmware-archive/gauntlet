package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"os/exec"
)

type Result struct {
	Pipeline  string `json:"pipeline"`
	Pipecount int    `json:"pipecount"`
	Stage     string  `json:"stage"`
	Stagecount int    `json:"stagecount"`
	Jobname   string  `json:"jobname"`
	Gitinfo   string  `json:"gitinfo"`
	Pass      bool    `json:"pass"`
}

func toInt(a string) int {
	val, err := strconv.Atoi(a)
	if err != nil {
		panic(fmt.Sprintf("could not convert '%s' to int: %s", a, err))
	}
	return val
}


func main() {

	r := Result{
		Pipeline:  os.Getenv("GO_PIPELINE_NAME"),
		Pipecount: toInt(os.Getenv("GO_PIPELINE_COUNTER")),
		Stage:  os.Getenv("GO_STAGE_NAME"),
		Stagecount: toInt(os.Getenv("GO_STAGE_COUNTER")),
		Jobname: os.Getenv("GO_JOB_NAME"),
		Gitinfo: os.Getenv("GO_REVISION"),
		Pass:    status,
	}

	json, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("json = %s\n", json)

	resp, err := http.Post("http://localhost:3000/results", "application/json", bytes.NewBuffer(json))

	body := bytes.NewBuffer(nil)
	io.Copy(body, resp.Body)

	fmt.Printf("posted: %s\nresp.Status is '%s'\nBody is '%s'\n", json, resp.Status, string(body.Bytes()))

	fmt.Printf("\n\n\n lets do a GET for comparison:\n")
	resp, err = http.Get("http://localhost:3000/results")
	if err != nil {
	   panic(err)
	}	

	fmt.Printf("client did GET, got response: %#v\n", resp)
}
