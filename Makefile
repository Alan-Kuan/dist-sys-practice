.PHONY: build clean

build:
	go build -o . ./cmd/...

clean:
	-rm `ls ./cmd`
