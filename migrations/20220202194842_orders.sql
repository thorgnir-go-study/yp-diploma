-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS gophermart.order_status
(
    id integer NOT NULL,
    name character varying NOT NULL,
    CONSTRAINT order_status_pkey PRIMARY KEY (id)
);

INSERT INTO gophermart.order_status (id, name) VALUES (1, 'NEW');
INSERT INTO gophermart.order_status (id, name) VALUES (2, 'PROCESSING');
INSERT INTO gophermart.order_status (id, name) VALUES (3, 'INVALID');
INSERT INTO gophermart.order_status (id, name) VALUES (4, 'PROCESSED');

CREATE TABLE IF NOT EXISTS gophermart."order"
(
    id uuid NOT NULL,
    order_number character varying NOT NULL,
    user_id uuid NOT NULL,
    status_id integer NOT NULL,
    accrual numeric,
    uploaded_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone,
    CONSTRAINT order_pkey PRIMARY KEY (id),
    CONSTRAINT order_order_status FOREIGN KEY (status_id)
        REFERENCES gophermart.order_status (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT order_user FOREIGN KEY (user_id)
        REFERENCES gophermart.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
