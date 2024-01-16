package reporter_test

import (
	"errors"

	reporterpkg "code.cloudfoundry.org/cpu-entitlement-admin-plugin/reporter"
	"code.cloudfoundry.org/cpu-entitlement-admin-plugin/reporter/reporterfakes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reporter", func() {
	var (
		reporter           reporterpkg.Reporter
		fakeCfClient       *reporterfakes.FakeCloudFoundryClient
		fakeMetricsFetcher *reporterfakes.FakeMetricsFetcher
	)

	BeforeEach(func() {
		fakeCfClient = new(reporterfakes.FakeCloudFoundryClient)
		fakeMetricsFetcher = new(reporterfakes.FakeMetricsFetcher)

		fakeCfClient.GetSpacesReturns([]reporterpkg.Space{
			{
				Name: "space1",
				Applications: []reporterpkg.Application{
					{Name: "app1", Guid: "space1-app1-guid"},
					{Name: "app2", Guid: "space1-app2-guid"},
				},
			},
			{
				Name: "space2",
				Applications: []reporterpkg.Application{
					{Name: "app1", Guid: "space2-app1-guid"},
				},
			},
		}, nil)

		fakeMetricsFetcher.FetchInstanceEntitlementUsagesStub = func(appGuid string) ([]float64, error) {
			switch appGuid {
			case "space1-app1-guid":
				return []float64{1.5, 0.5}, nil
			case "space1-app2-guid":
				return []float64{0.3}, nil
			case "space2-app1-guid":
				return []float64{0.2}, nil
			}

			return nil, nil
		}

		reporter = reporterpkg.New(fakeCfClient, fakeMetricsFetcher)
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
				fakeCfClient.GetSpacesReturns(nil, errors.New("get-space-error"))
			})

			It("returns the error", func() {
				Expect(err).To(MatchError("get-space-error"))
			})
		})

		When("getting the entitlement usage for an app fails", func() {
			BeforeEach(func() {
				fakeMetricsFetcher.FetchInstanceEntitlementUsagesReturns(nil, errors.New("fetch-error"))
			})

			It("returns the error", func() {
				Expect(err).To(MatchError("fetch-error"))
			})
		})
	})
})
