-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS gophermart."order"
ADD CONSTRAINT order_number_unique UNIQUE (order_number);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
