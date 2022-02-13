-- +goose Up
-- +goose StatementBegin
create table gophermart.processing_task
(
    id         uuid                      not null
        constraint processing_task_pk
            primary key,
    order_id   uuid                      not null
        constraint processing_task_order_id_fk
            references gophermart."order",
    to_run_at  timestamptz               not null,
    status     int                       not null,
    updated_at timestamptz default now() not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
