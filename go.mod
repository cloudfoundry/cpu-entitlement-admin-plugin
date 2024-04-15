module code.cloudfoundry.org/cpu-entitlement-admin-plugin

go 1.21.0

toolchain go1.21.8

// cf-cli still requires go-interact v1 so we don't want to bump that.
replace github.com/vito/go-interact => github.com/vito/go-interact v1.0.0 // indirect

require (
	code.cloudfoundry.org/cli v0.0.0-20220215194203-081f8aefeeef
	code.cloudfoundry.org/cpu-entitlement-plugin v0.0.0-20190820143339-d3a05702ac58
	code.cloudfoundry.org/go-log-cache/v2 v2.0.7
	github.com/google/uuid v1.6.0
	github.com/onsi/ginkgo/v2 v2.17.1
	github.com/onsi/gomega v1.32.0
)

require (
	code.cloudfoundry.org/bytefmt v0.0.0-20240409162623-53abdae41573 // indirect
	code.cloudfoundry.org/go-loggregator/v9 v9.2.0 // indirect
	code.cloudfoundry.org/jsonry v1.1.4 // indirect
	code.cloudfoundry.org/tlsconfig v0.0.0-20240410162701-78a97c114f7f // indirect
	github.com/SermoDigital/jose v0.9.2-0.20161205224733-f6df55f235c2 // indirect
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/bmatcuk/doublestar v1.3.4 // indirect
	github.com/charlievieth/fs v0.0.3 // indirect
	github.com/cloudfoundry/bosh-cli v6.4.1+incompatible // indirect
	github.com/cloudfoundry/bosh-utils v0.0.457 // indirect
	github.com/cppforlife/go-patch v0.2.0 // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/pprof v0.0.0-20240409012703-83162a5b38cd // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.1 // indirect
	github.com/jessevdk/go-flags v1.5.0 // indirect
	github.com/lunixbochs/vtclean v1.0.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/vito/go-interact v1.0.1 // indirect
	golang.org/x/crypto v0.22.0 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/term v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/tools v0.20.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240415141817-7cd4c1c1f9ec // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240415141817-7cd4c1c1f9ec // indirect
	google.golang.org/grpc v1.63.2 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/cheggaaa/pb.v1 v1.0.28 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
