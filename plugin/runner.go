package plugin

import "code.cloudfoundry.org/cpu-entitlement-admin-plugin/reporter"

type Reporter interface {
	OverEntitlementInstances() (reporter.Report, error)
}

type Renderer interface {
	Render(reporter.Report)
}

type Runner struct {
	reporter Reporter
	renderer Renderer
}

func NewRunner(reporter Reporter, renderer Renderer) *Runner {
	return &Runner{
		reporter: reporter,
		renderer: renderer,
	}
}

func (r *Runner) Run() error {
	report, _ := r.reporter.OverEntitlementInstances()

	r.renderer.Render(report)
	return nil
}
