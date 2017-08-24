
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE mls_listings (
    id serial PRIMARY KEY,
    mls_id varchar(128) NOT NULL,
    address varchar(512) NOT NULL,
    unit smallint NOT NULL,
    price smallint NOT NULL,

    description varchar(1024) DEFAULT '' NOT NULL,
    has_locker bool DEFAULT false NOT NULL,
    has_parking bool DEFAULT false NOT NULL,
    apt_size smallint DEFAULT 0 NOT NULL,
    exposure varchar(16) DEFAULT '' NOT NULL,
    distance jsonb,

    created_at timestamp DEFAULT now() NOT NULL,
    deleted_at timestamp,

    CONSTRAINT unique_mls_id UNIQUE (mls_id)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE mls_listings;
