#!/bin/bash

docker-compose run db psql -h db -U postgres <<-EOF
 $(cat migrations/1_orgs.sql)
EOF

docker-compose run db psql -h db -U postgres <<-EOF
 $(cat migrations/2_users.sql)
EOF

docker-compose run db psql -h db -U postgres <<-EOF
 $(cat migrations/3_user_orgs.sql)
EOF
