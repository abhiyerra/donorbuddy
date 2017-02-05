build: test
	go build

deps:
	go get -u
	go get -u github.com/mattes/migrate

test:
	go test -v

dev-run: build
	./donorbuddy config.dev.json

db-migrate:
	 migrate -url $$DONORBUDDY_DB -path ./migrations up

# http://apps.irs.gov/pub/epostcard/data-download-pub78.zip
# unzip -a data-download-pub78.zip
db-import:
	psql $$DONORBUDDY_DB -c "COPY orgs (ein, name, city, state, country, notes) FROM STDIN WITH DELIMITER AS '|'" < data-download-pub78.txt

release:
	cd frontend && npm run build
	kubectl set image donorbuddy=buildleft/donorbuddy:$(shell git rev-parse HEAD)
