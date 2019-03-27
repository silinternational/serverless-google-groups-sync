test:
	docker-compose run app ./codeship/test.sh

debug:
	docker-compose logs -t app

deploy:
	docker-compose run app sls deploy -v

remove:
	docker-compose run app sls remove