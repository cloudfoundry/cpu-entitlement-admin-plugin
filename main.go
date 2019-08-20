package main

import (
	"code.cloudfoundry.org/cli/plugin"
	cpuadminplugin "code.cloudfoundry.org/cpu-entitlement-admin-plugin/plugin"
)

func main() {
	plugin.Start(cpuadminplugin.New())
}
