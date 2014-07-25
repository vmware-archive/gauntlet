package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	cv "github.com/smartystreets/goconvey/convey"
)

func TestCaptureJobSuccessOrNotStatus(t *testing.T) {
	SetupFakeGoCdEnvVar()
	ws := NewWebServer("localhost:3000").Start()
	defer ws.Stop()
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

			ws := NewWebServer("localhost:3000").Start()
			defer ws.Stop()

			Watcher([]string{"/bin/echo", "hello", "gocd"})

			cv.So(ws.LastReqBody, cv.ShouldEqual, `{"pipeline":"jasons-fake-pipe","pipecount":"5","stage":"defaultStage","stagecount":"1","jobname":"defaultJob","gitinfo":"2e35c516660344a37cb5094102eb6d0e4f0414cc","pass":true}`)
		})
	})
}

func TestMissingEnvVariables(t *testing.T) {
	cv.Convey("Given missing GoCD env variables", t, func() {
		cv.Convey("watcher should record what GoCD env vars that are present, and not crash", func() {
			ClearFakeGoCdEnvVar()
			ws := NewWebServer("localhost:3000").Start()
			defer ws.Stop()

			Watcher([]string{"/bin/echo", "hello", "gocd"})

			cv.So(ws.LastReqBody, cv.ShouldEqual, `{"pipeline":"","pipecount":"","stage":"","stagecount":"","jobname":"","gitinfo":"","pass":true}`)

			Watcher([]string{"/bin/non-existant-binary-tried", "hello", "gocd"})

			cv.So(ws.LastReqBody, cv.ShouldEqual, `{"pipeline":"","pipecount":"","stage":"","stagecount":"","jobname":"","gitinfo":"","pass":false}`)

		})
	})
}

func TestGauntletEndpointEnvVariable(t *testing.T) {
	SetupFakeGoCdEnvVar()
	ws := NewWebServer("localhost:3000").Start()
	defer ws.Stop()

	cv.Convey("Given a missing or incorrect gauntlet endpoint env variables (GAUNTLET_HTTP_SERVER is not set)", t, func() {

		cv.Convey("for a missing endpoint, watcher should provide a helpful error and exit with a failure", func() {
			DeleteOneEnvVar("GAUNTLET_HTTP_SERVER")

			status, err := Watcher([]string{"/bin/echo", "hello", "gocd"})

			cv.So(status, cv.ShouldEqual, false)
			cv.So(strings.HasSuffix(err.Error(), "GAUNTLET_HTTP_SERVER is not set"), cv.ShouldEqual, true)

		})

		cv.Convey("for an incorrect endpoint, watcher should provide a helpful error and exit with a failure", func() {
			badTcpEnd := "localhost:3999"
			os.Setenv("GAUNTLET_HTTP_SERVER", badTcpEnd)
			// verify that the port is indeed not in use now.
			cv.So(PortIsBound(badTcpEnd), cv.ShouldEqual, false)

			status, err := Watcher([]string{"/bin/echo", "hello", "gocd"})

			cv.So(status, cv.ShouldEqual, false)
			expectedMsg := fmt.Sprintf("GAUNTLET_HTTP_SERVER cannot be contacted on '%s'", badTcpEnd)
			cv.So(strings.HasSuffix(err.Error(), expectedMsg), cv.ShouldEqual, true)

		})

	})
}

func SetupFakeGoCdEnvVar() {
	os.Setenv("GAUNTLET_HTTP_SERVER", "localhost:3000")
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

func DeleteOneEnvVar(key string) {
	env := os.Environ()
	m := EnvToMap(env)

	delete(m, key)

	os.Clearenv()
	MapToEnv(m)
}

func TestDeletingAnEnvVariable(t *testing.T) {
	cv.Convey("DeleteOneEnvVar removes an environmental variable from the environment", t, func() {
		os.Setenv("REMOVE_THIS", "NOTEMPTY")
		DeleteOneEnvVar("REMOVE_THIS")
		cv.So(os.Getenv("REMOVE_THIS"), cv.ShouldEqual, "")
	})
}
