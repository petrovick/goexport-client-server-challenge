down:
	docker-compose down

build:
	docker-compose build
	
up: down
	WATCH_FILES=0 docker-compose up -d

login:
	docker-compose run -w /usr/src/app go-client-server /bin/sh

logs:
	docker-compose logs --tail=10 -f

testing:
	docker-compose run -w /usr/src/app go-client-server /bin/sh -c "npm t"

testing-watch:
	docker-compose run -w /usr/src/app go-client-server /bin/sh -c "npm run test:watch"