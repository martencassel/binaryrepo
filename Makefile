.PHONY: test
test: 
	go test ./... -v

cover:
	go test ./... -cover
