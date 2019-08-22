.PHONY: test

help:
	@echo 'Help:'
	@echo '  build ........................ build the cpu entitlement admin binary'
	@echo '  install ...................... build and install the cpu entitlement admin binary'
	@echo '  test ......................... run tests (such as they are)'
	@echo '  help ......................... show help menu'

build:
	go build -mod vendor

test:
	ginkgo -r --race --skipPackage=e2e

install: build
	cf uninstall-plugin CPUEntitlementAdminPlugin || true
	cf install-plugin ./cpu-entitlement-admin-plugin -f
