build:
	@mkdir -p build
	go build -o build/moose .

run:
	go run .