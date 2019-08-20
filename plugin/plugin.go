package plugin // import "code.cloudfoundry.org/cpu-entitlement-admin-plugin/plugin"

import (
	"os"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/cf/trace"
	"code.cloudfoundry.org/cli/plugin"
	"code.cloudfoundry.org/cpu-entitlement-admin-plugin/output"
	"code.cloudfoundry.org/cpu-entitlement-admin-plugin/reporter"
)

type CPUEntitlementAdminPlugin struct{}

func New() CPUEntitlementAdminPlugin {
	return CPUEntitlementAdminPlugin{}
}

func (p CPUEntitlementAdminPlugin) Run(cli plugin.CliConnection, args []string) {
	if args[0] == "CLI-MESSAGE-UNINSTALL" {
		os.Exit(0)
	}

	traceLogger := trace.NewLogger(os.Stdout, true, os.Getenv("CF_TRACE"), "")
	ui := terminal.NewUI(os.Stdin, os.Stdout, terminal.NewTeePrinter(os.Stdout), traceLogger)

	reporter := reporter.New(cli)
	renderer := output.NewRenderer(ui)
	runner := NewRunner(reporter, renderer)

	err := runner.Run()
	if err != nil {
		ui.Failed(err.Error())
		os.Exit(1)
	}
}

func (p CPUEntitlementAdminPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "CPUEntitlementAdminPlugin",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 0,
			Build: 1,
		},
		Commands: []plugin.Command{
			{
				Name:     "over-entitlement-instances",
				Alias:    "oei",
				HelpText: "See which instances are over entitlement",
				UsageDetails: plugin.Usage{
					Usage: "cf over-entitlement-instances",
				},
			},
		},
	}
}
