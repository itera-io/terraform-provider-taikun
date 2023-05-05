TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=itera-io
NAMESPACE=dev
NAME=taikun
BINARY=terraform-provider-${NAME}
VERSION=0.1.0
OS_ARCH=darwin_amd64

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

clean:
	rm -f ${BINARY}

.PHONY: build dockerbuild commoninstall install dockerinstall test testacc clean
