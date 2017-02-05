CREATE TABLE user_orgs (
    user_id integer references users (id),
    org_id integer references orgs (id),
    created_at timestamp default now(),
    updated_at timestamp default now()
);
