TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=itera-io
NAMESPACE=dev
NAME=taikun
BINARY=terraform-provider-${NAME}
VERSION=0.1.0
OS_ARCH=linux_amd64

default: install

build:
	go build -o ${BINARY}

dockerbuild:
	DOCKER_BUILDKIT=1 docker build --rm --target bin --output . .

commoninstall:
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

install: build commoninstall

dockerinstall: dockerbuild commoninstall

test:
	go test -i $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

# --- Radek's rigorous testing here ---
#ACCEPTANCE_TESTS='(TestAccResourceTaikunRepository|TestAccDataSourceTaikunRepository|TestAccDataSourceTaikunCatalog|TestAccResourceTaikunCatalog|TestAccResourceTaikunOrganization|TestAccDataSourceTaikunOrganization|TestAccDataSourceTaikunShowback|TestAccDataSourceTaikunAccessProfile|TestAccResourceTaikunAccessProfile|TestAccDataSourceTaikunAlertingProfile|TestAccResourceTaikunProjectModifyAlertingProfile|TestAccResourceTaikunKubernetesProfile|TestAccResourceTaikunBilling|TestAccResourceTaikunProjectModifyImages|TestAccResourceTaikunProjectD|TestAccResourceTaikunProjectK|TestAccResourceTaikunUser|TestAccResourceTaikunShowback|TestAccResourceTaikunProjectModifyFlavors|TestAccResourceTaikunProjectU|TestAccResourceTaikunProject$$|TestAccDataSourceTaikunCloudCredentialOpenStack|TestAccDataSourceTaikunCloudCredentialsOpenStack|TestAccDataSourceTaikunCloudCredentialAzure|TestAccDataSourceTaikunCloudCredentialsAzure|TestAccDataSourceTaikunCloudCredentialAWS|TestAccDataSourceTaikunCloudCredentialsAWS|TestAccDataSourceTaikunKubernetesProfile|TestAccDataSourceTaikunBillingRule|TestAccDataSourceTaikunBillingCredential|TestAccDataSourceTaikunBackupCredential)' ./scripts/rerun_failed_tests.sh
# 1A
#ACCEPTANCE_TESTS='(TestAccDataSourceTaikunCloudCredentialsProxmoxWithFilter$$|TestAccDataSourceTaikunCloudCredentialsProxmox$$|TestAccResourceTaikunCloudCredentialProxmox$$|TestAccResourceTaikunCloudCredentialProxmoxLock$$|TestAccResourceTaikunCloudCredentialProxmoxUpdate$$|TestAccResourceTaikunAutoscalerOpenstackProject$$|TestAccResourceTaikunProjectWasm$$|TestAccDataSourceTaikunProject|TestAccDataSourceTaikunImagesDeprecated|TestAccDataSourceTaikunImagesOpenStack|TestAccDataSourceTaikunFlavorsOpenStack|TestAccDataSourceTaikunImagesAzure|TestAccDataSourceTaikunFlavorsAzure|TestAccDataSourceTaikunFlavorsAWS|TestAccDataSourceTaikunImagesAWS|TestAccResourceTaikunProjectE|TestAccDataSourceTaikunPolicyProfile|TestAccDataSourceTaikunSlackConfiguration|TestAccResourceTaikunStandaloneProfile|TestAccDataSourceTaikunStandaloneProfile|TestAccDataSourceTaikunUser|TestAccResourceTaikunProjectToggle|TestAccResourceTaikunSlack|TestAccResourceTaikunAlerting|TestAccResourceTaikunPolicyProfile|TestAccResourceTaikunCloudCredentials|TestAccResourceTaikunCloudCredentialOpenStack|TestAccResourceTaikunCloudCredentialAzure|TestAccResourceTaikunCloudCredentialAWS|TestAccResourceTaikunBackupCredential|TestProvider$$)' ./scripts/rerun_failed_tests.sh
# 1B
#ACCEPTANCE_TESTS='(TestAccResourceTaikunRepository|TestAccDataSourceTaikunRepository|TestAccDataSourceTaikunCatalog|TestAccResourceTaikunCatalog|TestAccResourceTaikunOrganization|TestAccDataSourceTaikunOrganization|TestAccDataSourceTaikunShowback|TestAccDataSourceTaikunAccessProfile|TestAccResourceTaikunAccessProfile|TestAccDataSourceTaikunAlertingProfile|TestAccResourceTaikunProjectModifyAlertingProfile|TestAccResourceTaikunKubernetesProfile|TestAccResourceTaikunBilling|TestAccResourceTaikunProjectModifyImages|TestAccResourceTaikunProjectD|TestAccResourceTaikunProjectK|TestAccResourceTaikunUser|TestAccResourceTaikunShowback|TestAccResourceTaikunProjectModifyFlavors|TestAccResourceTaikunProjectU|TestAccResourceTaikunProject$$|TestAccDataSourceTaikunCloudCredentialOpenStack|TestAccDataSourceTaikunCloudCredentialsOpenStack|TestAccDataSourceTaikunCloudCredentialAzure|TestAccDataSourceTaikunCloudCredentialsAzure|TestAccDataSourceTaikunCloudCredentialAWS|TestAccDataSourceTaikunCloudCredentialsAWS|TestAccDataSourceTaikunKubernetesProfile|TestAccDataSourceTaikunBillingRule|TestAccDataSourceTaikunBillingCredential|TestAccDataSourceTaikunBackupCredential)' ./scripts/rerun_failed_tests.sh
# 1C
#ACCEPTANCE_TESTS='(TestAccResourceTaikunCloudCredentialZadara|TestAccDataSourceTaikunCloudCredentialZadara|TestAccDataSourceTaikunCloudCredentialsZadara|TestAccDataSourceTaikunImagesZadara$$|TestAccResourceTaikunCloudCredentialGCP$$|TestAccDataSourceTaikunCloudCredentialsGCP$$|TestAccDataSourceTaikunImagesGCP$$|TestAccDataSourceTaikunCloudCredentialGCP$$)' ./scripts/rerun_failed_tests.sh

ACCEPTANCE_TESTS='(TestAccResourceTaikunProjectModifyFlavors$$)'

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


rtestacc:
	date
	go clean -testcache
	TF_ACC=1 go test . ./taikun/*/testing -v -run ${ACCEPTANCE_TESTS} -timeout 120m

rtestacc1:
	date
	# __________________________________________________________
	# >>>>>>>>>>>>>>>> TESTACC start Threads: 1 <<<<<<<<<<<<<<<<
	go clean -testcache
	TF_ACC=1 go test $(TEST) -v -run ${ACCEPTANCE_TESTS} -timeout 120m -parallel=1

rtestacc2:
	date
	# __________________________________________________________
	# >>>>>>>>>>>>>>>> TESTACC start Threads: 2 <<<<<<<<<<<<<<<<
	go clean -testcache
	TF_ACC=1 go test $(TEST) -v -run ${ACCEPTANCE_TESTS} -timeout 120m -parallel=2

rtestacc3:
	date
	# __________________________________________________________
	# >>>>>>>>>>>>>>>> TESTACC start Threads: 3 <<<<<<<<<<<<<<<<
	go clean -testcache
	TF_ACC=1 go test $(TEST) -v -run ${ACCEPTANCE_TESTS} -timeout 120m -parallel=3

rtestacc4:
	date
	# __________________________________________________________
	# >>>>>>>>>>>>>>>> TESTACC start Threads: 4 <<<<<<<<<<<<<<<<
	go clean -testcache
	TF_ACC=1 go test $(TEST) -v -run ${ACCEPTANCE_TESTS} -timeout 120m -parallel=4

rtestaccrigorous: rtestacc1 rtestacc2 rtestacc3 rtestacc4

clean:
	rm -f ${BINARY}

.PHONY: build dockerbuild commoninstall install dockerinstall test testacc clean


#TEST?=$$(go list ./... | grep -v 'vendor')
#HOSTNAME=itera-io
#NAMESPACE=dev
#NAME=taikun
#BINARY=terraform-provider-${NAME}
#VERSION=0.1.0
#OS_ARCH=linux_amd64
#
#default: install
#
#build:
#	go build -o ${BINARY}
#
#dockerbuild:
#	DOCKER_BUILDKIT=1 docker build --rm --target bin --output . .
#
#commoninstall:
#	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
#	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
#
#install: build commoninstall
#
#dockerinstall: dockerbuild commoninstall
#
#test:
#	go test -i "$(TEST)" || exit 1
#	echo "$(TEST)" | xargs -t -n4 go test "$(TESTARGS)" -timeout=30s -parallel=4
#
#testacc:
#	# __________________________________________________________
#	# ---------------- TESTACC start Threads: 1 ----------------
#	go clean -testcache
#	TF_ACC=1 go test $(TEST) -v "$TESTARGS" -timeout 120m
#
#testacc2:
#	clear
#	# __________________________________________________________
#	# >>>>>>>>>>>>>>>> TESTACC start Threads: 2 <<<<<<<<<<<<<<<<
#	go clean -testcache
#	TF_ACC=1 go test "$(TEST)" -v "$TESTARGS" -timeout 120m -parallel=2
#
#testacc3:
#	#
#	# __________________________________________________________
#	# >>>>>>>>>>>>>>>> TESTACC start Threads: 3 <<<<<<<<<<<<<<<<
#	go clean -testcache
#	TF_ACC=1 go test "$(TEST)" -v "$TESTARGS" -timeout 120m -parallel=3
#
#testacc4:
#	#
#	# __________________________________________________________
#	# >>>>>>>>>>>>>>>> TESTACC start Threads: 4 <<<<<<<<<<<<<<<<
#	go clean -testcache
#	TF_ACC=1 go test "$(TEST)" -v "$TESTARGS" -timeout 120m -parallel=4
#
#testrigorous: testacc2 testacc3 testacc4
#
#clean:
#	rm -f ${BINARY}
#
#.PHONY: build dockerbuild commoninstall install dockerinstall test testacc clean
