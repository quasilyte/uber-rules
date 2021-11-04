.PHONY: ci-tests

ci-tests:
	go test -v -race -count=2 .
