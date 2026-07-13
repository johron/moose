.PHONY: build release run clean

build:
	@mkdir -p build
	go build -o build/moose .

release:
	@mkdir -p build
	env CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o build/moose .

run:
	go run .

clean:
	rm -rf build/
