package reporter

import (
	"code.cloudfoundry.org/cli/plugin"
)

type Report struct {
	InstanceIDs []string
}

type Reporter struct{}

func New(cli plugin.CliConnection) *Reporter {
	return &Reporter{}
}

func (r *Reporter) OverEntitlementInstances() (Report, error) {
	return Report{InstanceIDs: []string{"one", "two", "three"}}, nil
}
