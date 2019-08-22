package plugin_test

import (
	"errors"

	"code.cloudfoundry.org/cpu-entitlement-admin-plugin/plugin"
	"code.cloudfoundry.org/cpu-entitlement-admin-plugin/plugin/pluginfakes"
	"code.cloudfoundry.org/cpu-entitlement-admin-plugin/reporter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Runner", func() {
	var (
		fakeReporter *pluginfakes.FakeReporter
		fakeRenderer *pluginfakes.FakeRenderer

		runner *plugin.Runner
		err    error
		report reporter.Report
	)

	BeforeEach(func() {
		fakeReporter = new(pluginfakes.FakeReporter)
		fakeRenderer = new(pluginfakes.FakeRenderer)

		report = reporter.Report{
			SpaceReports: []reporter.SpaceReport{
				{
					SpaceName: "space-1",
					Apps: []string{
						"app-1",
						"app-2",
					},
				}, {
					SpaceName: "space-2",
					Apps: []string{
						"app-1",
					},
				},
			},
		}

		fakeReporter.OverEntitlementInstancesReturns(report, nil)

		runner = plugin.NewRunner(fakeReporter, fakeRenderer)
	})

	JustBeforeEach(func() {
		err = runner.Run()
	})

	It("collects reports and renders them", func() {
		Expect(err).NotTo(HaveOccurred())
		Expect(fakeRenderer.RenderCallCount()).To(Equal(1))
		Expect(fakeRenderer.RenderArgsForCall(0)).To(Equal(report))
	})

	When("the reporter fails", func() {
		BeforeEach(func() {
			fakeReporter.OverEntitlementInstancesReturns(reporter.Report{}, errors.New("reporter-err"))
		})

		It("returns the error", func() {
			Expect(err).To(MatchError("reporter-err"))
		})
	})

	When("the renderer fails", func() {
		BeforeEach(func() {
			fakeRenderer.RenderReturns(errors.New("renderer-err"))
		})

		It("returns the error", func() {
			Expect(err).To(MatchError("renderer-err"))
		})
	})
})
