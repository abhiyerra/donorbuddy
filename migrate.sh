#!/bin/bash

set -e

docker-compose run db psql -h db -U postgres <<-EOF
 $(cat migrations/1_orgs.sql)
EOF
