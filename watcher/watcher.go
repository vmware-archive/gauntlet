package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Result struct {
	Pipeline   string `json:"pipeline"`
	Pipecount  string `json:"pipecount"`
	Stage      string `json:"stage"`
	Stagecount string `json:"stagecount"`
	Jobname    string `json:"jobname"`
	Gitinfo    string `json:"gitinfo"`
	Pass       bool   `json:"pass"`
}

func toInt(a string) int {
	val, err := strconv.Atoi(a)
	if err != nil {
		panic(fmt.Sprintf("could not convert '%s' to int: %s", a, err))
	}
	return val
}

func Watcher(command []string) (status bool, err error) {

	gauntletServer := os.Getenv("GAUNTLET_HTTP_SERVER")
	if gauntletServer == "" {
		return false, errors.New("watcher config error: GAUNTLET_HTTP_SERVER is not set")
	}

	if !PortIsBound(gauntletServer) {
		return false, fmt.Errorf("watcher config error: GAUNTLET_HTTP_SERVER cannot be contacted on '%s'", gauntletServer)
	}

	cmd := exec.Command(command[0], command[1:]...)
	err = cmd.Run()

	if err == nil {
		status = true
	}

	postToGauntlet(status, gauntletServer)

	return status, err
}

func postToGauntlet(status bool, gauntletServer string) {

	r := Result{
		Pipeline:   os.Getenv("GO_PIPELINE_NAME"),
		Pipecount:  os.Getenv("GO_PIPELINE_COUNTER"),
		Stage:      os.Getenv("GO_STAGE_NAME"),
		Stagecount: os.Getenv("GO_STAGE_COUNTER"),
		Jobname:    os.Getenv("GO_JOB_NAME"),
		Gitinfo:    os.Getenv("GO_REVISION"),
		Pass:       status,
	}

	json, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("json = %s\n", json)

	resp, err := http.Post("http://"+gauntletServer+"/results", "application/json", bytes.NewBuffer(json))
	if err != nil {
		panic(err)
	}

	body := bytes.NewBuffer(nil)
	io.Copy(body, resp.Body)

	fmt.Printf("posted: %s\nresp.Status is '%s'\nBody is '%s'\n", json, resp.Status, string(body.Bytes()))

}

func main() {
	command := os.Args[1:]
	if len(command) < 1 {
		log.Printf("watcher error: no command given. exiting 1.\n")
		log.Printf("   usage: watcher <command> <optional_arg1>...\n")
		log.Printf("   example: watcher run_test.sh arg1 arg2")
		os.Exit(1)
	}
	status, err := Watcher(command)
	if err != nil {
		log.Printf("watcher error: '%s'\n", err)
		os.Exit(1)
	}
	if status == false {
		cmdstr := strings.Join(command, " ")
		log.Printf("watcher error: executing my command line '%v' gave false status, exiting 1.\n", cmdstr)
		os.Exit(1)
	}
	os.Exit(0)
}
