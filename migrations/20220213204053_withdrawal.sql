-- +goose Up
-- +goose StatementBegin
create table gophermart.withdrawal
(
    id           uuid,
    order_number varchar                 not null,
    sum          numeric                 not null,
    processed_at timestamp with time zone NOT NULL DEFAULT now(),

    CONSTRAINT withdrawal_pk PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
