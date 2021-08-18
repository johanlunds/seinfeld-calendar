run: build
	docker run -it --rm --name my-running-app -p 3000:3000 my-golang-app

build:
	docker build -t my-golang-app .

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