run: build
	./seinfeld-calendar

build-docker:
	docker build -t my-golang-app .

# cross-compiled for Mac ARM64
build:
	docker run --rm -v "${PWD}":/usr/src/app -w /usr/src/app -e GOOS=darwin -e GOARCH=arm64 golang:1.16-alpine go build -v

run-docker:
	docker run -it --rm --name my-running-app my-golang-app

clean:
	rm ./seinfeld-calendar