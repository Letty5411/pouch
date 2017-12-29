package main

import (
	"github.com/alibaba/pouch/test/environment"

	"github.com/alibaba/pouch/test/request"
	"github.com/go-check/check"
	"net/url"
)

// APIContainerAttachSuite is the test suite for container attach API.
type APIContainerAttachSuite struct{}

func init() {
	check.Suite(&APIContainerAttachSuite{})
}

// SetUpTest does common setup in the beginning of each test.
func (suite *APIContainerAttachSuite) SetUpTest(c *check.C) {
	SkipIfFalse(c, environment.IsLinux)
}

// TestContainerAttachStdin tests attaching stdin is OK.
func (suite *APIContainerAttachSuite) TestContainerAttachStdin(c *check.C) {
	// TODO
	// path := "/containers/{name:.*}/attach"

	//cname := "TestContainerAttachStdin"
	//CreateBusyboxContainerOk(c, cname)
	//StartContainerOk(c, cname)

}

// TestContainerAttachNotFound
func (suite *APIContainerAttachSuite) TestContainerAttachNotFound(c *check.C) {
	cname := "TestContainerAttachNotFound"

	q := url.Values{}
	q.Set("stdin", "1")
	query := request.WithQuery(q)
	header := request.WithHeader("Content-Type", "text/plain")

	resp, err := request.Post("/containers/"+cname+"/attach", query, header)
	c.Assert(err, check.IsNil)
	
	// TODO: should return 404 once issue#470 is done
	CheckRespStatus(c, resp, 200)
}
