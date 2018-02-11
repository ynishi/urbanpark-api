.PHONY: build

build:
	goimports -w *.go
	GOOS=linux go build -o handler 
	zip handler.zip handler 
