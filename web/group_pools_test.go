// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/tsuru/tsuru-usage/repositories"
	"github.com/tsuru/tsuru/cmd/cmdtest"

	"gopkg.in/check.v1"
)

func (s *S) TestGroupPoolUsage(c *check.C) {
	groupData := `{
	"Name": "group 1",
	"Teams": ["team 1", "team 2"],
	"Pools": ["pool 1", "pool 2"]
}`
	usageData := `[
	{
		"Month": "January",
		"Usage": [
			{
				"Pool": "pool1",
				"Usage": 5
			},
			{
				"Pool": "pool2",
				"Usage": 7
			}
		]
	},
	{
		"Month": "February",
		"Usage": [
			{
				"Pool": "pool2",
				"Usage": 2
			}
		]
	}
]`
	repositories.Client.Transport = &cmdtest.Transport{Message: groupData, Status: http.StatusOK}
	Client.Transport = &cmdtest.Transport{Message: usageData, Status: http.StatusOK}
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/web/teamgroups/mygroup/pools/2017", nil)
	c.Assert(err, check.IsNil)
	m := runServer()
	c.Assert(m, check.NotNil)
	m.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	body := recorder.Body.String()
	c.Assert(strings.Contains(body, "January"), check.Equals, true)
	c.Assert(strings.Contains(body, "pool1"), check.Equals, true)
	c.Assert(strings.Contains(body, "5.00"), check.Equals, true)
	c.Assert(strings.Contains(body, "pool2"), check.Equals, true)
	c.Assert(strings.Contains(body, "7.00"), check.Equals, true)
	c.Assert(strings.Contains(body, "12.00"), check.Equals, true)
	c.Assert(strings.Contains(body, "February"), check.Equals, true)
	c.Assert(strings.Contains(body, "pool2"), check.Equals, true)
	c.Assert(strings.Contains(body, "2.00"), check.Equals, true)
	c.Assert(strings.Contains(body, "Total"), check.Equals, true)
	c.Assert(strings.Contains(body, "14.00"), check.Equals, true)
}

func (s *S) TestGroupPoolUsageAPIError(c *check.C) {
	repositories.Client.Transport = &cmdtest.Transport{Status: http.StatusInternalServerError}
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/web/teamgroups/mygroup/pools/2017", nil)
	c.Assert(err, check.IsNil)
	m := runServer()
	c.Assert(m, check.NotNil)
	m.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusInternalServerError)
}

func (s *S) TestGroupPoolUsageInvalidJSON(c *check.C) {
	repositories.Client.Transport = &cmdtest.Transport{Message: "invalid", Status: http.StatusOK}
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/web/teamgroups/mygroup/pools/2017", nil)
	c.Assert(err, check.IsNil)
	m := runServer()
	c.Assert(m, check.NotNil)
	m.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusInternalServerError)
}
