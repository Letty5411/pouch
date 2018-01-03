package main

import (
	"github.com/alibaba/pouch/test/environment"
	"github.com/alibaba/pouch/test/request"

	"github.com/go-check/check"
	"time"
)

// APIContainerExecStartSuite is the test suite for container exec start API.
type APIContainerExecStartSuite struct{}

func init() {
	check.Suite(&APIContainerExecStartSuite{})
}

// SetUpTest does common setup in the beginning of each test.
func (suite *APIContainerExecStartSuite) SetUpTest(c *check.C) {
	SkipIfFalse(c, environment.IsLinux)
}

// TestContainerExecStartOk tests start exec is OK.
func (suite *APIContainerExecStartSuite) TestContainerExecStartOk(c *check.C) {
	cname := "TestContainerCreateExecStartOk"

	CreateBusyboxContainerOk(c, cname)

	StartContainerOk(c, cname)

	got := CreateExecEchoOk(c, cname)
	c.Logf("got=%s",got)

	StartContainerExecOk(c, got, false, false)

	DelContainerForceOk(c, cname)
}

// TestContainerExecStartNotFound tests starting an non-existing execID return error.
func (suite *APIContainerExecStartSuite) TestContainerExecStartNotFound(c *check.C) {
	resp, err := request.Post("/exec/TestContainerExecStartNotFound/start")
	c.Assert(err, check.IsNil)
	CheckRespStatus(c, resp, 404)
}

// TestContainerExecStartStopped tests start a process in a stopped container return error.
func (suite *APIContainerExecStartSuite) TestContainerExecStartStopped(c *check.C) {
	cname := "TestContainerExecStartStopped"

	CreateBusyboxContainerOk(c, cname)

	StartContainerOk(c, cname)

	got := CreateExecEchoOk(c, cname)

	StopContainerOk(c, cname)

	resp, err := StartContainerExec(c, got, false, false)
	c.Assert(err, check.IsNil)
	CheckRespStatus(c, resp, 409)

	DelContainerForceOk(c, cname)
}

// TestContainerExecStartPaused tests start a process in a paused container return error.
func (suite *APIContainerExecStartSuite) TestContainerExecStartPaused(c *check.C) {
	cname := "TestContainerExecStartPaused"

	CreateBusyboxContainerOk(c, cname)

	StartContainerOk(c, cname)

	got := CreateExecEchoOk(c, cname)

	PauseContainerOk(c, cname)

	resp, err := StartContainerExec(c, got, false, false)
	c.Assert(err, check.IsNil)
	CheckRespStatus(c, resp, 409)

	UnpauseContainer(c, cname)

	StartContainerExecOk(c, got, false, false)

	DelContainerForceOk(c, cname)
}

// TestContainerExecStartDup tests start a process twice return error.
func (suite *APIContainerExecStartSuite) TestContainerExecStartDup(c *check.C) {
	cname := "TestContainerExecStartDup"

	CreateBusyboxContainerOk(c, cname)

	StartContainerOk(c, cname)

	got := CreateExecEchoOk(c, cname)

	StartContainerExecOk(c, got, false, false)

	// TODO: Add wait exec function when there is an inspect exec API
	time.Sleep(100 * time.Millisecond)

	resp, err := StartContainerExec(c, got, false, false)
	c.Assert(err, check.IsNil)
	CheckRespStatus(c, resp, 409)

	DelContainerForceOk(c, cname)
}
