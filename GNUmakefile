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

# Radek's acklowledgment testing
RADEK_TESTS='(TestAccResourceTaikunProjectStandaloneAWSMinimalUpdateFlavor$$)'

# All the tests NOT OK so far


# Part 1a - Lipton tests / Every other resource -- 231 s
#RADEK_TESTS='(TestAccResourceTaikunUser|TestAccResourceTaikunSlack|TestAccResourceTaikunShowback|TestAccResourceTaikunPolicyProfile|TestAccResourceTaikunKubernetesProfile|TestAccResourceTaikunCloudCredentials|TestAccResourceTaikunCloudCredentialOpenStack|TestAccResourceTaikunCloudCredentialAzure|TestAccResourceTaikunCloudCredentialAWS|TestAccResourceTaikunBilling|TestAccResourceTaikunAlerting|TestAccResourceTaikunAccess|TestAccResourceTaikunBackupCredential|TestProvider$$)'

# Part 1b - 7up tests / Non resource projects -- 230 s
#RADEK_TESTS='(TestAccResourceTaikunProject$$|TestAccResourceTaikunProjectE|TestAccResourceTaikunProjectD|TestAccResourceTaikunProjectK|TestAccResourceTaikunProjectU|TestAccResourceTaikunProjectToggle|TestAccResourceTaikunProjectModify)'

# Part 1c - Mirinda tests / Data sources fast -- 230 s
#RADEK_TESTS='(TestAccDataSourceTaikunPolicyProfile|TestAccDataSourceTaikunSlackConfiguration|TestAccResourceTaikunStandaloneProfile|TestAccDataSourceTaikunStandaloneProfile|TestAccDataSourceTaikunShowbackRule|TestAccDataSourceTaikunUser)'

# Part 1d - Mountain Dew tests / 246 s
#RADEK_TESTS='(TestAccDataSourceTaikunCloudCredentialOpenStack|TestAccDataSourceTaikunCloudCredentialsOpenStack|TestAccDataSourceTaikunCloudCredentialAzure|TestAccDataSourceTaikunCloudCredentialsAzure|TestAccDataSourceTaikunCloudCredentialAWS$$|TestAccDataSourceTaikunCloudCredentialsAWS|TestAccResourceTaikunAccessProfile|TestAccDataSourceTaikunKubernetesProfile|TestAccDataSourceTaikunBillingRule|TestAccDataSourceTaikunBillingCredential|TestAccDataSourceTaikunBackupCredential|TestAccDataSourceTaikunKubernetes)'

# Part 1e - Fanta tests / 150s s
#RADEK_TESTS='(TestAccDataSourceTaikunProject|TestAccResourceTaikunOrganization|TestAccDataSourceTaikunOrganization|TestAccDataSourceTaikunShowback|TestAccDataSourceTaikunAccessProfile|TestAccDataSourceTaikunAlertingProfile|TestAccDataSourceTaikunImagesDeprecated|TestAccDataSourceTaikunImagesOpenStack|TestAccDataSourceTaikunFlavorsOpenStack|TestAccDataSourceTaikunImagesAzure|TestAccDataSourceTaikunFlavorsAzure|TestAccDataSourceTaikunFlavorsAWS|TestAccDataSourceTaikunImagesAWS)'

# --- Creating resources, long ---
# Openstack has enough limit to run the 4 Standalone OpenStack in parallel
# Part 2a - Pepsi Max- 256 s
#RADEK_TESTS='(TestAccResourceTaikunProjectStandaloneOpenStackMinimal$$)'
# Part 2b - Pepsi Lime - 730 s
#RADEK_TESTS='(TestAccResourceTaikunProjectStandaloneOpenStackMinimalUpdateIP$$)'
# Part 2c - Pepsi Mango - 388 s
#RADEK_TESTS='(TestAccResourceTaikunProjectStandaloneOpenStackMinimalUpdateFlavor$$)'
# Part 2d - Pepsi Classic 230 s
#RADEK_TESTS='(TestAccResourceTaikunProjectStandaloneOpenStackMinimalWithVolumeType$$)'

# Part  5a - Monster energy 780 s - Creating k8s cluster in openstack - so Far openstack does not have resources to run it in parallel with other Openstack tests.
#RADEK_TESTS='(TestAccResourceTaikunProjectMinimal)'

# Part 3a - Kofola Classic 391 s
#RADEK_TESTS='(TestAccResourceTaikunProjectStandaloneAWSMinimal$$)'
# Part 3b - Kofola Lime 583 s
#RADEK_TESTS='(TestAccResourceTaikunProjectStandaloneAWSMinimalUpdateFlavor$$)'

# Part 4a - Coca Cola classic -- 515
#RADEK_TESTS='(TestAccResourceTaikunProjectStandaloneAzureMinimal$$)'
# Part 4b - Coca Cola zero -- 706 s
#RADEK_TESTS='(TestAccResourceTaikunProjectStandaloneAzureMinimalUpdateFlavor$$)'
# Part 4c - Coca Cola Cherry -- 464 s
#RADEK_TESTS='(TestAccResourceTaikunProjectStandaloneAzureMinimalWithVolumeType$$)'


rtestacc:
	date
	go clean -testcache
	TF_ACC=1 go test . ./taikun -v -run ${RADEK_TESTS} -timeout 120m

rtestacc1:
	date
	# __________________________________________________________
	# >>>>>>>>>>>>>>>> TESTACC start Threads: 1 <<<<<<<<<<<<<<<<
	go clean -testcache
	TF_ACC=1 go test $(TEST) -v -run ${RADEK_TESTS} -timeout 120m -parallel=1

rtestacc2:
	date
	# __________________________________________________________
	# >>>>>>>>>>>>>>>> TESTACC start Threads: 2 <<<<<<<<<<<<<<<<
	go clean -testcache
	TF_ACC=1 go test $(TEST) -v -run ${RADEK_TESTS} -timeout 120m -parallel=2

rtestacc3:
	date
	# __________________________________________________________
	# >>>>>>>>>>>>>>>> TESTACC start Threads: 3 <<<<<<<<<<<<<<<<
	go clean -testcache
	TF_ACC=1 go test $(TEST) -v -run ${RADEK_TESTS} -timeout 120m -parallel=3

rtestacc4:
	date
	# __________________________________________________________
	# >>>>>>>>>>>>>>>> TESTACC start Threads: 4 <<<<<<<<<<<<<<<<
	go clean -testcache
	TF_ACC=1 go test $(TEST) -v -run ${RADEK_TESTS} -timeout 120m -parallel=4

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
