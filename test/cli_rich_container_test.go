package main

import (
	"errors"
	"os"
	"runtime"
	"strings"

	"github.com/alibaba/pouch/test/command"
	"github.com/alibaba/pouch/test/environment"

	"github.com/go-check/check"
	"github.com/gotestyourself/gotestyourself/icmd"
	"fmt"
)

// PouchRichContainerSuite is the test suite fo rich container related CLI.
type PouchRichContainerSuite struct{}

func init() {
	check.Suite(&PouchRichContainerSuite{})
}

var centosImage string

// SetUpSuite does common setup in the beginning of each test suite.
func (suite *PouchRichContainerSuite) SetUpSuite(c *check.C) {
	SkipIfFalse(c, environment.IsLinux)
	SkipIfFalse(c, environment.IsRuncVersionSupportRichContianer)

	command.PouchRun("pull", busyboxImage).Assert(c, icmd.Success)

	// Use image from AliYun on AliOS.
	if environment.IsAliKernel() {
		centosImage = "reg.docker.alibaba-inc.com/alibase/alios7u2:latest"
	} else {
		centosImage = "registry.hub.docker.com/library/centos:latest"
	}
	command.PouchRun("pull", centosImage).Assert(c, icmd.Success)
}

//// TearDownSuite does common cleanup in the end of each test suite.
//func (suite *PouchRichContainerSuite) TearDownSuite(c *check.C) {
//	command.PouchRun("rmi", centosImage).Assert(c, icmd.Success)
//}

// isFileExistsInImage checks if the file exists in given image.
func isFileExistsInImage(image string, file string) (bool, error) {
	pc, _, _, _ := runtime.Caller(0)
	tmpname := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	var funcname string
	for i := range tmpname {
		funcname = tmpname[i]
	}

	if image == "" || file == "" {
		return false, errors.New("image is nil")
	}

	// check the existence of /sbin/init in image
	expect := icmd.Expected{
		ExitCode: 0,
		Out:      "Access",
	}
	err := command.PouchRun("run", "--name", funcname, image, "stat", file).Compare(expect)
	defer command.PouchRun("rm", "-f", funcname)

	return err == nil, nil
}

// checkPidofProcessIsOne checks the process of pid 1 is expected.
func checkPidofProcessIsOne(cname string, p string) bool {
	// Check the number 1 process is dumb-init
	cmd := "ps -ef |grep " + p + " |awk '{print $1}'"
	fmt.Printf("cmd=%s",cmd)
	expect := icmd.Expected{
		ExitCode: 0,
		Out:      "1",
	}
	err := command.PouchRun("exec", cname, "sh", "-c", cmd).Compare(expect)
	if err != nil {
		fmt.Printf("err=%s\n", err)
	}
	return err == nil
}

// checkPPid checks the ppid of process is expected.
func checkPPid(cname string, p string, ppid string) bool {
	// Check the number 1 process is dumb-init
	cmd := "ps -ef |grep " + p + " |awk '{print $3}'"
	expect := icmd.Expected{
		ExitCode: 0,
		Out:      ppid,
	}
	err := command.PouchRun("exec", cname, "sh", "-c", cmd).Compare(expect)
	return err == nil
}

// checkInitScriptWorks
func checkInitScriptWorks(c *check.C, cname string, image string, richmode string) {
	// Check run shell script works
	script := "/tmp/" + cname + ".sh"
	os.Remove(script)
	defer os.Remove(script)

	if _, err := os.Create("/tmp/" + cname + ".sh"); err != nil {
		c.Fatal("Fail to create %s file", script)
	}

	cmd := "echo touch /tmp/" + cname + " > " + script
	icmd.RunCommand("sh", "-c", cmd).Assert(c, icmd.Success)

	command.PouchRun("run", "-d", "--privileged", "-v", "/tmp:/tmp", "--rich", "--rich-mode", richmode, "--initscript",
		script, "--name", cname, busyboxImage, "sleep 10000").Assert(c, icmd.Success)

	// Check the number 1 process is dumb-init
	cmd = "stat /tmp" + cname
	expect := icmd.Expected{
		Out: "Access",
		Err: "",
	}
	err := command.PouchRun("exec", cname, "sh", "-c", cmd).Compare(expect)
	c.Assert(err, check.IsNil)

	command.PouchRun("rm", "-f", cname)

	// Check run a CMD works
}

// TestRichContainerDumbInitWorks check the dumb-init works.
func (suite *PouchRichContainerSuite) TestRichContainerDumbInitWorks(c *check.C) {
	SkipIfFalse(c, environment.IsDumbInitExist)
	pc, _, _, _ := runtime.Caller(0)
	tmpname := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	var funcname string
	for i := range tmpname {
		funcname = tmpname[i]
	}

	{
		command.PouchRun("run", "-d", "--rich", "--rich-mode", "dumb-init", "--name", funcname, busyboxImage, "sleep 10000").Assert(c, icmd.Success)

		c.Assert(checkPidofProcessIsOne(funcname, "dumb-init"), check.Equals, true)
		c.Assert(checkPPid(funcname, "sleep", "1"), check.Equals, true)

		command.PouchRun("rm", "-f", funcname)
	}

	{
		//checkInitScriptWorks(c, funcname, busyboxImage,"dumb-init" )
	}

}

// TestRichContainerWrongArgs check the wrong args of rich container.
func (suite *PouchRichContainerSuite) TestRichContainerDumbInitWrongArgs(c *check.C) {
	SkipIfFalse(c, environment.IsDumbInitExist)

	// TODO

	// Don't add '--rich' when use other rich container related options should fail.

}

// TestRichContainerSbinInitWorks check the initd works.
func (suite *PouchRichContainerSuite) TestRichContainerInitdWorks(c *check.C) {
	pc, _, _, _ := runtime.Caller(0)
	tmpname := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	var funcname string
	for i := range tmpname {
		funcname = tmpname[i]
	}

	ok, _ := isFileExistsInImage(centosImage, "/sbin/init")
	if !ok {
		c.Skip("/sbin/init doesn't exist in test image")
	}

	{
		// --privileged is MUST required
		command.PouchRun("run", "-d", "--privileged", "--rich", "--rich-mode", "sbin-init", "--name", funcname, centosImage, "sleep 10000").Assert(c, icmd.Success)

		c.Assert(checkPidofProcessIsOne(funcname, "/sbin/init"), check.Equals, true)
		c.Assert(checkPPid(funcname, "sleep", "1"), check.Equals, true)

		command.PouchRun("rm", "-f", funcname)
	}

	{
		//checkInitScriptWorks(c, funcname, centosImage,"sbin-init" )

	}
}

// TestRichContainerSystemdWorks check the systemd works.
func (suite *PouchRichContainerSuite) TestRichContainerSystemdWorks(c *check.C) {
	pc, _, _, _ := runtime.Caller(0)
	tmpname := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	var funcname string
	for i := range tmpname {
		funcname = tmpname[i]
	}

	ok, _ := isFileExistsInImage(centosImage, "/usr/lib/systemd/systemd")
	if !ok {
		c.Skip("/usr/lib/systemd/systemd doesn't exist in test image")
	}

	defer command.PouchRun("rm", "-f", funcname)

	{
		command.PouchRun("run", "-d", "--privileged", "--rich", "--rich-mode", "systemd", "--name", funcname, centosImage, "echo test").Assert(c, icmd.Success)

		c.Assert(checkPidofProcessIsOne(funcname, "/usr/lib/systemd/systemd"), check.Equals, true)
		c.Assert(checkPPid(funcname, "sleep", "1"), check.Equals, true)
	}

	{
		//checkInitScriptWorks(c, funcname, centosImage,"systemd" )
	}
}
