package main

import (
	"os"
	"testing"

	cv "github.com/smartystreets/goconvey/convey"
)

/*
 When I run a job,
    watcher errors with a helpful error message when called without a command
    watcher records the env variable settings at a http endpoint
    watcher handles missing env variables
    watcher allows configuration of the http endpoint

*/

func TestCaptureJobSuccessOrNotStatus(t *testing.T) {
	SetupFakeGoCdEnvVar()
	cv.Convey("Given a pipeline job command to be run", t, func() {
		cv.Convey("when the Watcher runs the job, success should be noted correctly", func() {

			status, err := Watcher([]string{"/bin/echo", "hello", "gocd"})
			cv.So(err, cv.ShouldEqual, nil)
			cv.So(status, cv.ShouldEqual, true)
		})

		cv.Convey("when the Watcher runs the job, not passing should be noted correctly", func() {

			status, err := Watcher([]string{"/bin/echo-does-not-exist", "hello", "gocd"})
			cv.So(err, cv.ShouldNotEqual, nil)
			cv.So(status, cv.ShouldEqual, false)

		})
	})
}

func TestRecordEnvVariablesToHttpEndpoint(t *testing.T) {
	SetupFakeGoCdEnvVar()
	cv.Convey("Given a pipeline job command to be run", t, func() {
		cv.Convey("after watcher runs the job, watcher should record the GoCD env vars to an http endpoint", func() {

			addr := "localhost:3000"
			webserv := NewWebServer(addr)
			webserv.Start()
			defer webserv.Stop()

			Watcher([]string{"/bin/echo", "hello", "gocd"})

			lastReq := webserv.LastReqBody
			cv.So(lastReq, cv.ShouldNotEqual, "")
			cv.So(lastReq, cv.ShouldEqual, `{"pipeline":"jasons-fake-pipe","pipecount":"5","stage":"defaultStage","stagecount":"1","jobname":"defaultJob","gitinfo":"2e35c516660344a37cb5094102eb6d0e4f0414cc","pass":true}`)
		})
	})
}

func TestMissingEnvVariables(t *testing.T) {
	cv.Convey("Given missing GoCD env variables", t, func() {
		cv.Convey("watcher should record what GoCD env vars that are present, and not crash", func() {
			ClearFakeGoCdEnvVar()

			addr := "localhost:3000"
			webserv := NewWebServer(addr)
			webserv.Start()
			defer webserv.Stop()

			Watcher([]string{"/bin/echo", "hello", "gocd"})

			lastReq := webserv.LastReqBody
			cv.So(lastReq, cv.ShouldNotEqual, "")
			cv.So(lastReq, cv.ShouldEqual, `{"pipeline":"","pipecount":"","stage":"","stagecount":"","jobname":"","gitinfo":"","pass":true}`)

			Watcher([]string{"/bin/non-existant-binary-tried", "hello", "gocd"})

			lastReq = webserv.LastReqBody
			cv.So(lastReq, cv.ShouldNotEqual, "")
			cv.So(lastReq, cv.ShouldEqual, `{"pipeline":"","pipecount":"","stage":"","stagecount":"","jobname":"","gitinfo":"","pass":false}`)

		})
	})
}

func SetupFakeGoCdEnvVar() {
	os.Setenv("GO_SERVER_URL", "https://10.10.48.5:8154/go/")
	os.Setenv("GO_TRIGGER_USER", "releng")
	os.Setenv("GO_PIPELINE_NAME", "jasons-fake-pipe")
	os.Setenv("GO_PIPELINE_COUNTER", "5")
	os.Setenv("GO_PIPELINE_LABEL", "5")
	os.Setenv("GO_STAGE_NAME", "defaultStage")
	os.Setenv("GO_STAGE_COUNTER", "1")
	os.Setenv("GO_JOB_NAME", "defaultJob")
	os.Setenv("GO_REVISION", "2e35c516660344a37cb5094102eb6d0e4f0414cc")
	os.Setenv("GO_TO_REVISION", "2e35c516660344a37cb5094102eb6d0e4f0414cc")
	os.Setenv("GO_FROM_REVISION", "2e35c516660344a37cb5094102eb6d0e4f0414cc")
}

func ClearFakeGoCdEnvVar() {
	env := os.Environ()
	m := EnvToMap(env)

	delete(m, "GO_SERVER_URL")
	delete(m, "GO_TRIGGER_USER")
	delete(m, "GO_PIPELINE_NAME")
	delete(m, "GO_PIPELINE_COUNTER")
	delete(m, "GO_PIPELINE_LABEL")
	delete(m, "GO_STAGE_NAME")
	delete(m, "GO_STAGE_COUNTER")
	delete(m, "GO_JOB_NAME")
	delete(m, "GO_REVISION")
	delete(m, "GO_TO_REVISION")
	delete(m, "GO_FROM_REVISION")

	os.Clearenv()
	MapToEnv(m)
}
