

debug:
	docker-compose up -d app
	docker-compose exec app bash

test:
	docker-compose run app ./codeship/test.sh

slsdeploy:
	docker-compose up -d app
	docker-compose run app bash -c "sls deploy -v"