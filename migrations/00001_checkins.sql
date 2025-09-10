-- +goose Up
CREATE TABLE IF NOT EXISTS checkins (
    id SERIAL PRIMARY KEY,
    lat INT NOT NULL,
    lon INT NOT NULL,
    checkin_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
-- +goose Down
DROP TABLE IF EXISTS checkins;