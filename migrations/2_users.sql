CREATE TABLE users (
    id serial primary key,
    facebook_id varchar(255),
    stripe_customer_id varchar(255),
    stripe_subscription_id varchar(255),
    created_at timestamp default now(),
    updated_at timestamp default now()
);

ALTER TABLE users ADD CONSTRAINT users_facebook_id UNIQUE (facebook_id);
