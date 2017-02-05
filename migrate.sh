#!/bin/bash

psql -h $DONORBUDDY_HOST -U $DONORBUDDY_USER -W <<-EOF
 $(cat migrations/1_orgs.sql)
EOF

psql -h $DONORBUDDY_HOST -U $DONORBUDDY_USER -W <<-EOF
 $(cat migrations/2_users.sql)
EOF

psql -h $DONORBUDDY_HOST -U $DONORBUDDY_USER -W <<-EOF
 $(cat migrations/3_user_orgs.sql)
EOF
