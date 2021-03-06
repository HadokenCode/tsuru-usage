// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plan

import (
	"testing"

	"github.com/tsuru/tsuru-usage/db"

	"os"

	check "gopkg.in/check.v1"
)

var _ = check.Suite(&S{})

type S struct {
	conn *db.Storage
}

func Test(t *testing.T) { check.TestingT(t) }

func (s *S) SetUpTest(c *check.C) {
	var err error
	os.Setenv("MONGODB_DATABASE_NAME", "tsuru_usage_plan_tests")
	s.conn, err = db.Conn()
	c.Assert(err, check.IsNil)
	s.conn.PlanCosts().Database.DropDatabase()
}

func (s *S) TearDownSuite(c *check.C) {
	conn, err := db.Conn()
	c.Assert(err, check.IsNil)
	err = conn.TeamGroups().Database.DropDatabase()
	c.Assert(err, check.IsNil)
}

func (s *S) TestListCosts(c *check.C) {
	plan1 := PlanCost{
		Service:     "service",
		Plan:        "small",
		Type:        ServicePlan,
		Cost:        0.5,
		MeasureUnit: "dollars",
	}
	_, err := Save(plan1)
	c.Assert(err, check.IsNil)
	plan2 := PlanCost{
		Plan:        "small",
		Type:        AppPlan,
		Cost:        1,
		MeasureUnit: "dollars",
	}
	_, err = Save(plan2)
	c.Assert(err, check.IsNil)
	plans, err := ListCosts()
	c.Assert(err, check.IsNil)
	c.Assert(plans, check.DeepEquals, []PlanCost{plan1, plan2})
}

func (s *S) TestListServicesCosts(c *check.C) {
	plan1 := PlanCost{
		Service:     "service",
		Plan:        "small",
		Type:        ServicePlan,
		Cost:        0.5,
		MeasureUnit: "dollars",
	}
	_, err := Save(plan1)
	c.Assert(err, check.IsNil)
	plan2 := PlanCost{
		Plan:        "small",
		Type:        AppPlan,
		Cost:        1,
		MeasureUnit: "dollars",
	}
	_, err = Save(plan2)
	c.Assert(err, check.IsNil)
	plans, err := ListServicesCosts()
	c.Assert(err, check.IsNil)
	c.Assert(plans, check.DeepEquals, []PlanCost{plan1})
}

func (s *S) TestListAppsCosts(c *check.C) {
	plan1 := PlanCost{
		Service:     "service",
		Plan:        "small",
		Type:        ServicePlan,
		Cost:        0.5,
		MeasureUnit: "dollars",
	}
	_, err := Save(plan1)
	c.Assert(err, check.IsNil)
	plan2 := PlanCost{
		Plan:        "small",
		Type:        AppPlan,
		Cost:        1,
		MeasureUnit: "dollars",
	}
	_, err = Save(plan2)
	c.Assert(err, check.IsNil)
	plans, err := ListAppsCosts()
	c.Assert(err, check.IsNil)
	c.Assert(plans, check.DeepEquals, []PlanCost{plan2})
}

func (s *S) TestSave(c *check.C) {
	plan1 := PlanCost{
		Service:     "service",
		Plan:        "small",
		Type:        ServicePlan,
		Cost:        0.5,
		MeasureUnit: "dollars",
	}
	created, err := Save(plan1)
	c.Assert(err, check.IsNil)
	c.Assert(created, check.Equals, true)
	plans, err := ListCosts()
	c.Assert(err, check.IsNil)
	c.Assert(plans, check.DeepEquals, []PlanCost{plan1})
	plan1 = PlanCost{
		Service:     "service",
		Plan:        "small",
		Type:        ServicePlan,
		Cost:        1,
		MeasureUnit: "dollars",
	}
	created, err = Save(plan1)
	c.Assert(created, check.Equals, false)
	plans, err = ListCosts()
	c.Assert(err, check.IsNil)
	c.Assert(plans, check.DeepEquals, []PlanCost{plan1})
}
