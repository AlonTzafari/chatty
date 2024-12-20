-- +goose Up
-- +goose StatementBegin
CREATE TABLE messages (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    channel_id uuid REFERENCES channels (id) ON UPDATE CASCADE ON DELETE CASCADE,
    user_id uuid REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE,
    content TEXT NOT NULL,
    createdAt TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE messages;
-- +goose StatementEnd
