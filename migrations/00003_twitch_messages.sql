-- +goose Up
CREATE TABLE twitch_messages (
    id SERIAL PRIMARY KEY,
    username VARCHAR(25) NOT NULL,
    display_name VARCHAR(25),
    message TEXT NOT NULL,
    color VARCHAR(7),
    badges JSONB,
    is_mod BOOLEAN DEFAULT FALSE,
    is_vip BOOLEAN DEFAULT FALSE,
    is_broadcaster BOOLEAN DEFAULT FALSE,
    time TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    INDEX idx_twitch_messages_time (time),
    INDEX idx_twitch_messages_username (username)
);

-- +goose Down
DROP TABLE twitch_messages;