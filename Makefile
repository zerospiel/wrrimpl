BLD_FLAGS=CGO_ENABLED=1
LOCAL_BIN=$(CURDIR)/bin

.PHONY: .build
.build:
	@$(BLD_FLAGS) go build -o $(LOCAL_BIN)/wrrimpl ./cmd/wrrimpl/main.go

.PHONY: build
build: .build

.PHONY: .test
.test:
	@cd $(CURDIR)/wrr && go test -race -count=1 -timeout=30s ./...

.PHONY: test
test: .test
