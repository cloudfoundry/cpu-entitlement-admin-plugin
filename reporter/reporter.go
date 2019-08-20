package reporter // import "code.cloudfoundry.org/cpu-entitlement-admin-plugin/reporter"

import (
	"context"
	"fmt"

	"code.cloudfoundry.org/cli/plugin"
	plugin_models "code.cloudfoundry.org/cli/plugin/models"
	logcache "code.cloudfoundry.org/log-cache/pkg/client"
)

type Report struct {
	SpaceReports []SpaceReport
}

type SpaceReport struct {
	SpaceName string
	Apps      []string
}

type Reporter struct {
	cli            plugin.CliConnection
	logCacheClient *logcache.Client
}

func New(cli plugin.CliConnection, logCacheClient *logcache.Client) Reporter {
	return Reporter{
		cli:            cli,
		logCacheClient: logCacheClient,
	}
}

func (r Reporter) OverEntitlementInstances() (Report, error) {
	spaceReports := []SpaceReport{}

	spaces, _ := r.cli.GetSpaces()
	for _, space := range spaces {
		spaceModel, err := r.cli.GetSpace(space.Name)
		if err != nil {
			return Report{}, err
		}

		apps, err := r.getApps(spaceModel)
		if err != nil {
			return Report{}, err
		}

		spaceReports = append(spaceReports, SpaceReport{SpaceName: space.Name, Apps: apps})
	}

	return Report{SpaceReports: spaceReports}, nil
}

func (r Reporter) getApps(spaceModel plugin_models.GetSpace_Model) ([]string, error) {
	apps := []string{}
	for _, app := range spaceModel.Applications {
		isOverEntitlement, err := r.isOverEntitlement(app.Guid)
		if err != nil {
			return nil, err
		}
		if isOverEntitlement {
			apps = append(apps, app.Name)
		}
	}
	return apps, nil
}

func (r Reporter) isOverEntitlement(appGuid string) (bool, error) {
	appInstancesUsages, err := r.logCacheClient.PromQL(context.Background(), fmt.Sprintf(`absolute_usage{source_id="%s"} / absolute_entitlement{source_id="%s"}`, appGuid, appGuid))
	if err != nil {
		return false, err
	}

	isOverEntitlement := false
	for _, sample := range appInstancesUsages.GetVector().GetSamples() {
		if sample.GetPoint().GetValue() > 1 {
			isOverEntitlement = true
		}
	}

	return isOverEntitlement, nil
}
