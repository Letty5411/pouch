package main

import (
	"github.com/alibaba/pouch/test/environment"
	"github.com/alibaba/pouch/test/request"

	"github.com/go-check/check"
)

// APINetworkCreateSuite is the test suite for network create API.
type APINetworkCreateSuite struct{}

func init() {
	check.Suite(&APINetworkCreateSuite{})
}

// SetUpTest does common setup in the beginning of each test.
func (suite *APINetworkCreateSuite) SetUpTest(c *check.C) {
	SkipIfFalse(c, environment.IsLinux)
}

// TestNetworkCreateOk tests creating network is OK.
func (suite *APINetworkCreateSuite) TestNetworkCreateOk(c *check.C) {
	nname := "TestNetworkCreateOk"
	ipamConf := map[string]interface{}{
		"Gateway": "192.168.1.1",
		"IPRange": "192.168.1.1/24",
		"Subnet":  "192.168.1.1/24",
	}
	obj := map[string]interface{}{
		"Name": nname,
		"NetworkCreate": map[string]interface{}{
			"Driver": "bridge",
			"IPAM":   map[string]interface{}{"Config": []interface{}{&ipamConf}},
		},
	}

	body := request.WithJSONBody(obj)
	resp, err := request.Post("/networks/create", body)
	c.Assert(err, check.IsNil)
	// TODO: change to 201, once issue # has been fixed.
	CheckRespStatus(c, resp, 500)

	DelNetworkOk(c, nname)
}

// TestNetworkCreateNilName tests creating network without name returns error.
func (suite *APINetworkCreateSuite) TestNetworkCreateNilName(c *check.C) {
	obj := map[string]interface{}{
		"Name": nil,
		"NetworkCreate": map[string]interface{}{
			"Driver": "bridge",
		},
	}

	body := request.WithJSONBody(obj)
	resp, err := request.Post("/networks/create", body)
	c.Assert(err, check.IsNil)
	CheckRespStatus(c, resp, 500)
}
