down:
	docker-compose down

build:
	docker-compose build
	
up: down
	docker-compose up -d

login:
	docker-compose run -w /usr/src/app go-client /bin/sh

logs:
	docker-compose logs --tail=10 -f