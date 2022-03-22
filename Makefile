.PHONY: cmd

dev: | build serve

serve:
	./streamcast

build:
	go build .

cmd:
	go build -o debug cmd/debug.go
	go build -o copy cmd/copy.go
