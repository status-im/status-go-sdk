lint-install:
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install

lint:
	@echo "lint"
	@gometalinter ./...

UNIT_TEST_PACKAGES := $(shell go list ./...)

test:
	go test -coverpkg= $(UNIT_TEST_PACKAGES)

dev-deps:
	go get -u github.com/stretchr/testify
	go get -u github.com/golang/mock/gomock
	go get -u github.com/golang/mock/mockgen

mock:
	mockgen -package=sdk  -destination=sdk_mock.go -source=sdk.go
