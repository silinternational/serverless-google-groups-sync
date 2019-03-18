

debug:
	docker-compose up -d app
	docker-compose exec app bash

test:
	docker-compose run app ./codeship/test.sh

slsdeploy:
	docker-compose up -d app
	docker-compose exec app bash -c "./run-slsdeploy.sh"