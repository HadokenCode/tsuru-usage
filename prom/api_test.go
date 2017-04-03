package prom

import (
	"testing"
	"time"

	check "gopkg.in/check.v1"

	"golang.org/x/net/context"

	"github.com/prometheus/client_golang/api/prometheus"
	"github.com/prometheus/common/model"
)

type S struct{}

func Test(t *testing.T) { check.TestingT(t) }

type fakeQueryAPI struct {
	q string
	t time.Time
	m model.Value
	e error
}

func (q *fakeQueryAPI) Query(ctx context.Context, query string, ts time.Time) (model.Value, error) {
	return q.m, q.e
}

func (q *fakeQueryAPI) QueryRange(ctx context.Context, query string, r prometheus.Range) (model.Value, error) {
	return nil, nil
}

func (s *S) TestAvgOverPeriod(c *check.C) {
	api := &fakeQueryAPI{
		m: model.Vector{
			{Value: model.SampleValue(10)},
		},
	}
	Client = &prometheusAPI{
		queryAPI: api,
	}
	d := time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC)
	f, err := GetAvgOverPeriod("metric{label=\"a\"}", "10d", d)
	c.Assert(err, check.IsNil)
	c.Assert(f, check.Equals, 10)
	c.Assert(api.q, check.DeepEquals, "avg(avg_over_time(metric{label=\"a\"}[10d]))")
	c.Assert(api.t, check.DeepEquals, d)
}