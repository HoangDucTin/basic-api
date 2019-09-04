
export GOPATH=/home/auttasakl/Project/go

PROJECT=go-api

check: build
	@./${PROJECT};

build:
	@go build;

clean:
	@rm -f ${PROJECT}

# End-of-file

