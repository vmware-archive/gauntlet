package main

import (
	"testing"
	cv "github.com/smartystreets/goconvey/convey"
)

/*
 When I run a job,
    shep captures the success/fail status in pass
    shep records the env variable settings at a http endpoint
    shep handles missing env variables
    shep allows configuration of the http endpoint

*/

func TestCaptureJobSuccessOrFailStatus(t *testing.T) {
	cv.Convey("Given a pipeline job command to be run", t, func() {
		cv.Convey("when shep runs the job, success should be noted correctly", func() {

		  status, err := Shep("/bin/echo", []string{"hello", "gocd"})
		  cv.So(status, cv.ShouldEqual, true)
		})

		cv.Convey("when shep runs the job, failure should be noted correctly", func() {

		  status, err := Shep("/bin/echo-does-not-exist", []string{"hello", "gocd"})
		  cv.So(status, cv.ShouldEqual, false)

		})
	})
}

func TestRecordEnvVariablesToHttpEndpoint(t *testing.T) {
	cv.Convey("Given a pipeline job command to be run", t, func() {
		cv.Convey("after shep runs the job, shep should record the GoCD env vars to an http endpoint", func() {
		})
	})
}

func TestMissingEnvVariables(t *testing.T) {
	cv.Convey("Given missing GoCD env variables", t, func() {
		cv.Convey("shep should record what GoCD env vars that are present, and not crash", func() {
		})
	})
}


func TestHttpEndpointConfig(t *testing.T) {
	cv.Convey("shep should allow configuration of http endpoint", t, func() {

	})
}
