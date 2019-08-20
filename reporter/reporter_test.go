package reporter_test

import (
	"fmt"

	plugin_models "code.cloudfoundry.org/cli/plugin/models"
	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	reporterpkg "code.cloudfoundry.org/cpu-entitlement-admin-plugin/reporter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reporter", func() {
	var (
		reporter          reporterpkg.Reporter
		fakeCliConnection *pluginfakes.FakeCliConnection
	)

	BeforeEach(func() {
		fakeCliConnection = new(pluginfakes.FakeCliConnection)
		fakeCliConnection.GetSpacesReturns([]plugin_models.GetSpaces_Model{
			{Guid: "space1-guid", Name: "space1"},
			{Guid: "space2-guid", Name: "space2"},
			{Guid: "space3-guid", Name: "space3"},
		}, nil)

		fakeCliConnection.GetSpaceStub = func(spaceName string) (plugin_models.GetSpace_Model, error) {
			switch spaceName {
			case "space1":
				return plugin_models.GetSpace_Model{
					Applications: []plugin_models.GetSpace_Apps{
						{Name: "app1", Guid: "space1-app1-guid"},
						{Name: "app2", Guid: "space1-app2-guid"},
						{Name: "app3", Guid: "space1-app3-guid"},
					},
				}, nil
			case "space2":
				return plugin_models.GetSpace_Model{
					Applications: []plugin_models.GetSpace_Apps{
						{Name: "app1", Guid: "space2-app1-guid"},
					},
				}, nil
			case "space3":
				return plugin_models.GetSpace_Model{
					Applications: []plugin_models.GetSpace_Apps{
						{Name: "app1", Guid: "space3-app1-guid"},
					},
				}, nil
			}

			return plugin_models.GetSpace_Model{}, fmt.Errorf("Space '%s' not found", spaceName)
		}

		reporter = reporterpkg.New(fakeCliConnection)
	})

	Describe("OverEntitlementInstances", func() {
		It("returns all instances that are over entitlement", func() {
			report, err := reporter.OverEntitlementInstances()
			Expect(err).NotTo(HaveOccurred())

			Expect(report).To(Equal(reporterpkg.Report{
				SpaceReports: []reporterpkg.SpaceReport{
					reporterpkg.SpaceReport{
						SpaceName: "space1",
						Apps: []string{
							"app1",
							"app2",
						},
					},
					reporterpkg.SpaceReport{
						SpaceName: "space2",
						Apps: []string{
							"app1",
						},
					},
				},
			}))
		})
	})
})
