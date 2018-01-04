package main

import (
	"github.com/alibaba/pouch/test/command"
	"github.com/alibaba/pouch/test/environment"
	"github.com/go-check/check"
	"github.com/gotestyourself/gotestyourself/icmd"
	"strings"
)

// CRIRunPodSandboxSuite is the test suite for CRI RunPodSandboxSuite interface.
type CRIRunPodSandboxSuite struct{}

func init() {
	check.Suite(&CRIRunPodSandboxSuite{})
}

// SetUpTest does common setup in the beginning of each test.
func (suite *CRIRunPodSandboxSuite) SetUpTest(c *check.C) {
	SkipIfFalse(c, environment.IsLinux)
}

// TestRunPodSandboxSuiteWorks tests RunPodSandboxSuite could work.
func (suite *CRIRunPodSandboxSuite) TestRunPodSandboxSuiteWorks(c *check.C) {
	// TODO
	command.PouchRun("cri", "info").Assert(c, icmd.Success)
}
