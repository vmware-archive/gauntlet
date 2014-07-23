package main

import (
	"fmt"
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

func TestCaptureJobSuccessOrFailStatus(t *testing.T) {
	cv.Convey("Given a pipeline job command to be run", t, func() {
		cv.Convey("when the Watcher runs the job, success should be noted correctly", func() {

			status, err := Watcher([]string{"/bin/echo", "hello", "gocd"})
			cv.So(err, cv.ShouldEqual, nil)
			cv.So(status, cv.ShouldEqual, true)
		})

		cv.Convey("when the Watcher runs the job, failure should be noted correctly", func() {

			status, err := Watcher([]string{"/bin/echo-does-not-exist", "hello", "gocd"})
			cv.So(err, cv.ShouldNotEqual, nil)
			cv.So(status, cv.ShouldEqual, false)

		})
	})
}

func TestRecordEnvVariablesToHttpEndpoint(t *testing.T) {
	cv.Convey("Given a pipeline job command to be run", t, func() {
		cv.Convey("after watcher runs the job, watcher should record the GoCD env vars to an http endpoint", func() {

			StartWebServer()

			Watcher([]string{"/bin/echo", "hello", "gocd"})

			fmt.Printf("%#v\n", LastServerRequest)
			cv.So(LastServerRequest, cv.ShouldEqual, false)
		})
	})
}

func TestMissingEnvVariables(t *testing.T) {
	cv.Convey("Given missing GoCD env variables", t, func() {
		cv.Convey("watcher should record what GoCD env vars that are present, and not crash", func() {
		})
	})
}

func TestHttpEndpointConfig(t *testing.T) {
	cv.Convey("watcher should allow configuration of http endpoint", t, func() {

	})
}
