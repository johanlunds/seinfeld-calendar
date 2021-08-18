run: build
	docker run -it --rm --name my-running-app -p 3000:3000 my-golang-app

run-mac: build-mac
	./seinfeld-calendar

build:
	docker build -t my-golang-app .

# cross-compiled for Mac ARM64
build-mac:
	docker run --rm -v "${PWD}":/usr/src/app -w /usr/src/app -e GOOS=darwin -e GOARCH=arm64 golang:1.16-alpine go build -v

run-shell:
	docker run -it --rm --name my-running-app my-golang-app sh

# Run "heroku login" first
deploy:
	DOCKER_DEFAULT_PLATFORM=linux/amd64 heroku container:push --app seinfeldcalendar web
	heroku container:release --app seinfeldcalendar web

open:
	open https://seinfeldcalendar.herokuapp.com/

clean:
	rm ./seinfeld-calendar