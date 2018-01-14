OUTPUT_DIR=./_output
SOURCE_DIRS = .
PACKAGES := go list ./... | grep -v /vendor | grep -v /out

bin=candid
product: clean
#	env GOOS=linux GOARCH=amd64 go build -o ${OUTPUT_DIR}/${bin}.linux
	go build -o ${OUTPUT_DIR}/${bin}.linux

build: clean
	go build -o ${OUTPUT_DIR}/${bin}

test: clean
	@go test -v -race ./pkg/...

.PHONY: fmtcheck
fmtcheck:
	@gofmt -s -l $(SOURCE_DIRS) | grep ".*\.go"; if [ "$$?" = "0" ]; then exit 1; fi
	
.PHONY: vet
vet:
	@go vet $(shell $(PACKAGES))

clean:
	@: if [ -f ${OUTPUT_DIR} ] then rm -rf ${OUTPUT_DIR} fi
