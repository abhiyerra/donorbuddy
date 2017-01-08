#!/bin/bash

set -e

docker-compose run db psql -h db -U postgres -c "COPY orgs (ein, name, city, state, country, notes) FROM STDIN WITH DELIMITER AS '|'" < data-download-pub78.txt
