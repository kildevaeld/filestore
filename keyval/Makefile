
BINARY="keyval"
VERSION?="0.0.1"
OUTDIR?=${PWD}

build:
	@echo "Running go get"
	@go get  ./...
	@echo "Running go build"
	@go build -ldflags "-s -X github.com/kildevaeld/torsten/torsten/cmd.VERSION=${VERSION}" -o ${OUTDIR}/${BINARY}
