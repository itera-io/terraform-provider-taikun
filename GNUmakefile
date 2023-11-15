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

# TF_ACC=1 go test github.com/itera-io/terraform-provider-taikun github.com/itera-io/terraform-provider-taikun/taikun -v -run TestAccResourceTaikunBillingCredential$ -timeout=30s -parallel=4
# TF_ACC=1 go test . ./taikun -v -run TestAccResourceTaikunBillingCredential$ -timeout=30s -parallel=4
# TESTARGS="-run (TestAccResourceTaikunBillingCredential$|TestAccResourceTaikunBillingRule$)" make rtestacc
# TESTARGS="-run (TestAccResourceTaikunBillingCredential$|TestAccResourceTaikunBillingRule$)" make rtestaccrigorous
# RARGS="-run (TestAccResourceTaikunBillingCredential$|TestAccResourceTaikunBillingRule$)" make rtestacc

# TF_ACC=1 go test github.com/itera-io/terraform-provider-taikun github.com/itera-io/terraform-provider-taikun/taikun -v -run TestAccResourceTaikunBillingCredential$ -timeout=30s -parallel=4

# Radek's acklowledgment testing
#RADEK_TESTS='(TestAccResourceTaikunProject$$)'

# All the tests NOT OK so far


# All the tests OK so far together
#RADEK_TESTS='(TestAccResourceTaikunProject$$|TestAccResourceTaikunProjectE|TestAccResourceTaikunProjectD|TestAccResourceTaikunProjectK|TestAccResourceTaikunProjectU|TestAccResourceTaikunProjectModify|TestAccResourceTaikunProjectToggle|TestAccResourceTaikunUser|TestAccResourceTaikunStandaloneProfile|TestAccResourceTaikunSlack|TestAccResourceTaikunShowback|TestAccResourceTaikunPolicyProfile|TestAccResourceTaikunOrganization|TestAccResourceTaikunKubernetesProfile|TestAccResourceTaikunCloudCredentials|TestAccResourceTaikunCloudCredentialOpenStack|TestAccResourceTaikunCloudCredentialAzure|TestAccResourceTaikunCloudCredentialAWS|TestAccResourceTaikunBilling|TestAccResourceTaikunAlerting|TestAccResourceTaikunAccess|TestAccResourceTaikunBackupCredential)'

# Part 1 - Every other resource
#RADEK_TESTS='(TestAccResourceTaikunUser|TestAccResourceTaikunStandaloneProfile|TestAccResourceTaikunSlack|TestAccResourceTaikunShowback|TestAccResourceTaikunPolicyProfile|TestAccResourceTaikunOrganization|TestAccResourceTaikunKubernetesProfile|TestAccResourceTaikunCloudCredentials|TestAccResourceTaikunCloudCredentialOpenStack|TestAccResourceTaikunCloudCredentialAzure|TestAccResourceTaikunCloudCredentialAWS|TestAccResourceTaikunBilling|TestAccResourceTaikunAlerting|TestAccResourceTaikunAccess|TestAccResourceTaikunBackupCredential)'
# Part 2 - Non resource projects
RADEK_TESTS='(TestAccResourceTaikunProject$$|TestAccResourceTaikunProjectE|TestAccResourceTaikunProjectD|TestAccResourceTaikunProjectK|TestAccResourceTaikunProjectU|TestAccResourceTaikunProjectToggle|TestAccResourceTaikunProjectModify)'
# Part 3 - Resource projects
# Part 4 - All data sources

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
