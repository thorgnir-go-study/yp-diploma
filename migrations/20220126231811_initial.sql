-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS gophermart;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS gophermart.users
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    login character varying  NOT NULL,
    password character varying NOT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT login_unique UNIQUE (login)
)

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
