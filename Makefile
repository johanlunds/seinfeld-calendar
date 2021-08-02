build:
	docker build -t my-golang-app .

build-mounted:
	docker run --rm -v "${PWD}":/usr/src/app -w /usr/src/app golang:1.16-alpine go build -v

build-crosscompiled:
	docker run --rm -v "${PWD}":/usr/src/app -w /usr/src/app -e GOOS=darwin -e GOARCH=arm64 golang:1.16-alpine go build -v

run:
	docker run -it --rm --name my-running-app my-golang-app

clean:
	rm ./seinfeld-calendar