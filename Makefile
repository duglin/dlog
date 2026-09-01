all: dotest

dotest:
	@go clean -testcache
	@go test -failfast -v ./test | sed "s/^        //"
