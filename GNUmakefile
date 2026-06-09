GOLANGCI_LINTERS_VERSION := v2.12.2

TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=itera-io
NAMESPACE=dev
NAME=taikun
BINARY=terraform-provider-${NAME}
VERSION=0.1.0
OS_ARCH=linux_amd64

## Including environmental variables necessary for running unit tests
include ./dev-env.sh

default: help

deps: go-linters-install check-terraform ## Installing development prerequisites locally and checking all dependencies (if they're installed)

check-terraform:
	@command -v terraform >/dev/null 2>&1 || { echo >&2 "Error: Terraform is not installed. Aborting."; exit 1; }
	@echo "Terraform is installed!"

build: go-vendor ## Builds Golang binary
	go build -o build/_output/${BINARY}

generate: build ## Generates Terraform's bindings
	go generate ./...

go-linters-install: ## Installs Golang's linters locally for verification
	curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(shell go env GOPATH)/bin ${GOLANGCI_LINTERS_VERSION}

go-vet:
	go vet ./...

lint: go-vet go-linters-install ## Runs golangci-lint against codebase
	golangci-lint run --timeout 5m

dockerbuild: ## Builds Docker image
	DOCKER_BUILDKIT=1 docker build --rm --target bin --output . .

commoninstall: ## Installs built binary to the host's system
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv build/_output/${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

install: ## Builds and installs Terraform provider locally
install: build commoninstall

dockerinstall: ## Builds Docker image and installs binary to the Host system
dockerinstall: dockerbuild commoninstall

test: ## Runs unit tests
	go test -i $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: ## Runs unit tests with specified arguments
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

go-tidy: ## Runs go mod tidy
	go mod tidy

go-vendor: go-tidy ## Runs go mod tidy && go mod vendor
	go mod vendor

clean-vendor: ## Removes vendor folder
	rm -rf vendor

# --- Radek's rigorous testing here ---
#ACCEPTANCE_TESTS='(TestAccResourceTaikunRepository|TestAccDataSourceTaikunRepository|TestAccDataSourceTaikunCatalog|TestAccResourceTaikunCatalog|TestAccResourceTaikunOrganization|TestAccDataSourceTaikunOrganization|TestAccDataSourceTaikunShowback|TestAccDataSourceTaikunAccessProfile|TestAccResourceTaikunAccessProfile|TestAccDataSourceTaikunAlertingProfile|TestAccResourceTaikunProjectModifyAlertingProfile|TestAccResourceTaikunKubernetesProfile|TestAccResourceTaikunBilling|TestAccResourceTaikunProjectModifyImages|TestAccResourceTaikunProjectD|TestAccResourceTaikunProjectK|TestAccResourceTaikunUser|TestAccResourceTaikunShowback|TestAccResourceTaikunProjectModifyFlavors|TestAccResourceTaikunProjectU|TestAccResourceTaikunProject$$|TestAccDataSourceTaikunCloudCredentialOpenStack|TestAccDataSourceTaikunCloudCredentialsOpenStack|TestAccDataSourceTaikunCloudCredentialAzure|TestAccDataSourceTaikunCloudCredentialsAzure|TestAccDataSourceTaikunCloudCredentialAWS|TestAccDataSourceTaikunCloudCredentialsAWS|TestAccDataSourceTaikunKubernetesProfile|TestAccDataSourceTaikunBillingRule|TestAccDataSourceTaikunBillingCredential|TestAccDataSourceTaikunBackupCredential)' ./scripts/rerun_failed_tests.sh
# 1A
#ACCEPTANCE_TESTS='(TestAccDataSourceTaikunCloudCredentialsProxmoxWithFilter$$|TestAccDataSourceTaikunCloudCredentialsProxmox$$|TestAccResourceTaikunCloudCredentialProxmox$$|TestAccResourceTaikunCloudCredentialProxmoxLock$$|TestAccResourceTaikunCloudCredentialProxmoxUpdate$$|TestAccResourceTaikunAutoscalerOpenstackProject$$|TestAccResourceTaikunProjectWasm$$|TestAccDataSourceTaikunProject|TestAccDataSourceTaikunImagesDeprecated|TestAccDataSourceTaikunImagesOpenStack|TestAccDataSourceTaikunFlavorsOpenStack|TestAccDataSourceTaikunImagesAzure|TestAccDataSourceTaikunFlavorsAzure|TestAccDataSourceTaikunFlavorsAWS|TestAccDataSourceTaikunImagesAWS|TestAccResourceTaikunProjectE|TestAccDataSourceTaikunPolicyProfile|TestAccDataSourceTaikunSlackConfiguration|TestAccResourceTaikunStandaloneProfile|TestAccDataSourceTaikunStandaloneProfile|TestAccDataSourceTaikunUser|TestAccResourceTaikunProjectToggle|TestAccResourceTaikunSlack|TestAccResourceTaikunAlerting|TestAccResourceTaikunPolicyProfile|TestAccResourceTaikunCloudCredentials|TestAccResourceTaikunCloudCredentialOpenStack|TestAccResourceTaikunCloudCredentialAzure|TestAccResourceTaikunCloudCredentialAWS|TestAccResourceTaikunBackupCredential|TestProvider$$)' ./scripts/rerun_failed_tests.sh
# 1B
#ACCEPTANCE_TESTS='(TestAccResourceTaikunRepository|TestAccDataSourceTaikunRepository|TestAccDataSourceTaikunCatalog|TestAccResourceTaikunCatalog|TestAccResourceTaikunOrganization|TestAccDataSourceTaikunOrganization|TestAccDataSourceTaikunShowback|TestAccDataSourceTaikunAccessProfile|TestAccResourceTaikunAccessProfile|TestAccDataSourceTaikunAlertingProfile|TestAccResourceTaikunProjectModifyAlertingProfile|TestAccResourceTaikunKubernetesProfile|TestAccResourceTaikunBilling|TestAccResourceTaikunProjectModifyImages|TestAccResourceTaikunProjectD|TestAccResourceTaikunProjectK|TestAccResourceTaikunUser|TestAccResourceTaikunShowback|TestAccResourceTaikunProjectModifyFlavors|TestAccResourceTaikunProjectU|TestAccResourceTaikunProject$$|TestAccDataSourceTaikunCloudCredentialOpenStack|TestAccDataSourceTaikunCloudCredentialsOpenStack|TestAccDataSourceTaikunCloudCredentialAzure|TestAccDataSourceTaikunCloudCredentialsAzure|TestAccDataSourceTaikunCloudCredentialAWS|TestAccDataSourceTaikunCloudCredentialsAWS|TestAccDataSourceTaikunKubernetesProfile|TestAccDataSourceTaikunBillingRule|TestAccDataSourceTaikunBillingCredential|TestAccDataSourceTaikunBackupCredential)' ./scripts/rerun_failed_tests.sh
# 1C
#ACCEPTANCE_TESTS='(TestAccResourceTaikunCloudCredentialZadara|TestAccDataSourceTaikunCloudCredentialZadara|TestAccDataSourceTaikunCloudCredentialsZadara|TestAccDataSourceTaikunImagesZadara$$|TestAccResourceTaikunCloudCredentialGCP$$|TestAccDataSourceTaikunCloudCredentialsGCP$$|TestAccDataSourceTaikunImagesGCP$$|TestAccDataSourceTaikunCloudCredentialGCP$$)' ./scripts/rerun_failed_tests.sh

ACCEPTANCE_TESTS='(TestAcc.*TaikunOrganization.*)'

# --- CI: Not creating resources ---
# Acklowledgment testing ALPHA
# Acklowledgment testing BETA

# --- CI: Creating resources ---
# Openstack user (found in Bitwarden) has enough limit to run the 4 Standalone OpenStack in parallel
# Part 2a - Openstack - 256 s
# Part 2b - Openstack - 730 s
# Part 2c - Openstack - 388 s
# Part 2d - Openstack - 230 s

# Part  5a - Openstack - cca 780 s
# Creating k8s cluster in openstack - so Far openstack does not have resources to run it in parallel with other Openstack tests.

# Part 3a - AWS 1 391 s (this one is with VM spots)
# Part 3b - AWS 2 583 s

# Part 4a - Azure -- 515
# Part 4b - Azure -- 706 s
# Part 4c - Azure -- 464 s

rtestacc: install
	date
	go clean -testcache
	TF_ACC=1 go test . ./taikun/*/testing -v -run ${ACCEPTANCE_TESTS} -timeout 120m

rtestacc1: install
	date
	# __________________________________________________________
	# >>>>>>>>>>>>>>>> TESTACC start Threads: 1 <<<<<<<<<<<<<<<<
	go clean -testcache
	TF_ACC=1 go test $(TEST) -v -run ${ACCEPTANCE_TESTS} -timeout 120m -parallel=1

rtestacc2: install
	date
	# __________________________________________________________
	# >>>>>>>>>>>>>>>> TESTACC start Threads: 2 <<<<<<<<<<<<<<<<
	go clean -testcache
	TF_ACC=1 go test $(TEST) -v -run ${ACCEPTANCE_TESTS} -timeout 120m -parallel=2

rtestacc3: install
	date
	# __________________________________________________________
	# >>>>>>>>>>>>>>>> TESTACC start Threads: 3 <<<<<<<<<<<<<<<<
	go clean -testcache
	TF_ACC=1 go test $(TEST) -v -run ${ACCEPTANCE_TESTS} -timeout 120m -parallel=3

rtestacc4: install
	date
	# __________________________________________________________
	# >>>>>>>>>>>>>>>> TESTACC start Threads: 4 <<<<<<<<<<<<<<<<
	go clean -testcache
	TF_ACC=1 go test $(TEST) -v -run ${ACCEPTANCE_TESTS} -timeout 120m -parallel=4

rtestaccrigorous: rtestacc1 rtestacc2 rtestacc3 rtestacc4

clean: clean-vendor ## Removes built binary
	rm -f build/_output/${BINARY}

test-ci: ## Simulates CI/CD pipeline locally
test-ci: build dockerbuild commoninstall install dockerinstall test testacc clean

.PHONY: help
help: # Credits to https://gist.github.com/prwhite/8168133 for this handy oneliner
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
