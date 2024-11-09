-- +goose Up
-- +goose StatementBegin
CREATE TABLE channels (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    avatar VARCHAR(255),
    createdAt TIMESTAMP NOT NULL
);

CREATE TABLE channels_users (
    channel_id  uuid REFERENCES channels (id) ON UPDATE CASCADE ON DELETE CASCADE,
    user_id     uuid REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT channels_users_pkey PRIMARY KEY (channel_id, user_id)
);

ALTER TABLE users
ADD avatar varchar(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN avatar;

DROP TABLE channels_users;

DROP TABLE channels;
-- +goose StatementEnd
