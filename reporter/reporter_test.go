package reporter_test

import (
	"context"
	"errors"
	"fmt"
	"strings"

	plugin_models "code.cloudfoundry.org/cli/plugin/models"
	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	reporterpkg "code.cloudfoundry.org/cpu-entitlement-admin-plugin/reporter"
	"code.cloudfoundry.org/cpu-entitlement-admin-plugin/reporter/reporterfakes"
	"code.cloudfoundry.org/log-cache/pkg/client"
	"code.cloudfoundry.org/log-cache/pkg/rpc/logcache_v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reporter", func() {
	var (
		reporter           reporterpkg.Reporter
		fakeCliConnection  *pluginfakes.FakeCliConnection
		fakeLogCacheClient *reporterfakes.FakeLogCacheClient
	)

	BeforeEach(func() {
		fakeCliConnection = new(pluginfakes.FakeCliConnection)
		fakeLogCacheClient = new(reporterfakes.FakeLogCacheClient)

		fakeCliConnection.GetSpacesReturns([]plugin_models.GetSpaces_Model{
			{Guid: "space1-guid", Name: "space1"},
			{Guid: "space2-guid", Name: "space2"},
		}, nil)

		fakeCliConnection.GetSpaceStub = func(spaceName string) (plugin_models.GetSpace_Model, error) {
			switch spaceName {
			case "space1":
				return plugin_models.GetSpace_Model{
					Applications: []plugin_models.GetSpace_Apps{
						{Name: "app1", Guid: "space1-app1-guid"},
						{Name: "app2", Guid: "space1-app2-guid"},
					},
				}, nil
			case "space2":
				return plugin_models.GetSpace_Model{
					Applications: []plugin_models.GetSpace_Apps{
						{Name: "app1", Guid: "space2-app1-guid"},
					},
				}, nil
			}

			return plugin_models.GetSpace_Model{}, fmt.Errorf("Space '%s' not found", spaceName)
		}

		fakeLogCacheClient.PromQLStub = func(_ context.Context, query string, _ ...client.PromQLOption) (*logcache_v1.PromQL_InstantQueryResult, error) {
			switch {
			case strings.Contains(query, "space1-app1-guid"):
				return instantQueryResult(
					sample("0", 1.5),
					sample("1", 0.5),
				), nil
			case strings.Contains(query, "space1-app2-guid"):
				return instantQueryResult(
					sample("0", 0.3),
				), nil
			case strings.Contains(query, "space1-app1-guid"):
				return instantQueryResult(
					sample("0", 0.2),
				), nil
			}

			return instantQueryResult(), nil
		}

		reporter = reporterpkg.New(fakeCliConnection, fakeLogCacheClient)
	})

	Describe("OverEntitlementInstances", func() {
		var (
			report reporterpkg.Report
			err    error
		)

		JustBeforeEach(func() {
			report, err = reporter.OverEntitlementInstances()
		})

		It("succeeds", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns all instances that are over entitlement", func() {
			Expect(report).To(Equal(reporterpkg.Report{
				SpaceReports: []reporterpkg.SpaceReport{
					reporterpkg.SpaceReport{
						SpaceName: "space1",
						Apps: []string{
							"app1",
						},
					},
				},
			}))
		})

		When("fetching the list of apps fails", func() {
			BeforeEach(func() {
				fakeCliConnection.GetSpaceReturns(plugin_models.GetSpace_Model{}, errors.New("get-space-error"))
			})

			It("returns the error", func() {
				Expect(err).To(MatchError("get-space-error"))
			})
		})

		When("getting the entitlement usage for an app fails", func() {
			BeforeEach(func() {
				fakeLogCacheClient.PromQLReturns(nil, errors.New("promql-error"))
			})

			It("returns the error", func() {
				Expect(err).To(MatchError("promql-error"))
			})
		})
	})
})

func instantQueryResult(samples ...*logcache_v1.PromQL_Sample) *logcache_v1.PromQL_InstantQueryResult {
	return &logcache_v1.PromQL_InstantQueryResult{
		Result: &logcache_v1.PromQL_InstantQueryResult_Vector{
			Vector: &logcache_v1.PromQL_Vector{
				Samples: samples,
			},
		},
	}
}

func sample(time string, value float64) *logcache_v1.PromQL_Sample {
	return &logcache_v1.PromQL_Sample{
		Point: &logcache_v1.PromQL_Point{
			Time:  time,
			Value: value,
		},
	}
}
