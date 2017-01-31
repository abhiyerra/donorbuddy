build: test
	go build

deps:
	docker-compose pull
	go get -u

test:
	#go test -v

dev-run: build
	docker-compose build
	docker-compose up

dev-clean:
	docker-compose stop
	docker-compose rm

dev-db-migrate:
	 bash ./migrate.sh

dev-db-import:
	# http://apps.irs.gov/pub/epostcard/data-download-pub78.zip
	# unzip -a data-download-pub78.zip
	bash ./import-orgs.sh


release:
	cd frontend && npm run build
