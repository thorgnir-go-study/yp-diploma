-- +goose Up
-- +goose StatementBegin
alter table gophermart.withdrawal
    add user_id uuid not null;

alter table gophermart.withdrawal
    add constraint withdrawal_users_id_fk
        foreign key (user_id) references gophermart.users;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
