CREATE TABLE orgs (
    id serial primary key,
    name varchar(255),
    ein varchar(30),
    address varchar(255),
    city varchar(100),
    state varchar(20),
    country varchar(100),
    category varchar(50),
    verified bool default false,
    stripe_customer_id varchar(255),
    notes varchar(255),
    created_at timestamp default now(),
    updated_at timestamp default now()
);

ALTER TABLE orgs ADD CONSTRAINT orgs_ein UNIQUE (ein);

CREATE INDEX orgs_name_idx ON orgs USING GIN (to_tsvector('english', name));
