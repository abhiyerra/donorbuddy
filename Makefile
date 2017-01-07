build: test
	go build

deps:
	docker-compose pull
	go get -u

test:
	#go test -v

dev-run: build
	docker-compose up

dev-clean:
	docker-compose stop
	docker-compose rm

dev-db-migrate:
	bash ./migrate.sh
