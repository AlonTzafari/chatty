-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id uuid primary key,
    username VARCHAR(255),
    password VARCHAR(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
