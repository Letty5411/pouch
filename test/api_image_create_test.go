package main

import (
	"net/url"

	"github.com/alibaba/pouch/test/environment"
	"github.com/alibaba/pouch/test/request"

	"github.com/go-check/check"
)

// APIImageCreateSuite is the test suite for image create API.
type APIImageCreateSuite struct{}

func init() {
	check.Suite(&APIImageCreateSuite{})
}

// SetUpTest does common setup in the beginning of each test.
func (suite *APIImageCreateSuite) SetUpTest(c *check.C) {
	SkipIfFalse(c, environment.IsLinux)
}

// TestImageCreateOk tests creating an image is OK.
func (suite *APIImageCreateSuite) TestImageCreateOk(c *check.C) {
	q := url.Values{}
	q.Add("fromImage", helloworldImage)
	q.Add("tag", "latest")
	path := "/images/create"
	query := request.WithQuery(q)
	resp, err := request.Post(path, query)
	c.Assert(err, check.IsNil)
	c.Assert(resp.StatusCode, check.Equals, 200)

	resp, err = request.Delete("/images/" + helloworldImage)
	c.Assert(err, check.IsNil)
	c.Assert(resp.StatusCode, check.Equals, 204)
}

// TestImageCreateNil tests fromImage is nil.
func (suite *APIImageCreateSuite) TestImageCreateNil(c *check.C) {
	q := url.Values{}
	q.Add("fromImage", "")

	path := "/images/create"
	query := request.WithQuery(q)

	resp, err := request.Post(path, query)
	c.Assert(err, check.IsNil)
	c.Assert(resp.StatusCode, check.Equals, 400)
}

// TestImageCreateWithoutTag tests creating an image without tag, will use "latest" by default.
func (suite *APIImageCreateSuite) TestImageCreateWithoutTag(c *check.C) {
	q := url.Values{}
	q.Add("fromImage", helloworldImage)
	path := "/images/create"
	query := request.WithQuery(q)
	resp, err := request.Post(path, query)
	c.Assert(err, check.IsNil)
	c.Assert(resp.StatusCode, check.Equals, 200)

	resp, err = request.Delete("/images/" + helloworldImage)
	c.Assert(err, check.IsNil)
	c.Assert(resp.StatusCode, check.Equals, 204)
}

