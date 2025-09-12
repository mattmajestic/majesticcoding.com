-- +goose Up

-- YouTube stats history
CREATE TABLE IF NOT EXISTS youtube_stats (
    id SERIAL PRIMARY KEY,
    channel_id VARCHAR(255),
    subscriber_count INT,
    video_count INT,
    view_count BIGINT,
    recorded_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- GitHub stats history  
CREATE TABLE IF NOT EXISTS github_stats (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255),
    public_repos INT,
    followers INT,
    following INT,
    total_stars INT,
    recorded_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Twitch stats history
CREATE TABLE IF NOT EXISTS twitch_stats (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255),
    follower_count INT,
    view_count INT,
    is_live BOOLEAN,
    recorded_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- LeetCode stats history
CREATE TABLE IF NOT EXISTS leetcode_stats (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255),
    solved_count INT,
    ranking INT,
    main_language VARCHAR(100),
    recorded_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS youtube_stats;
DROP TABLE IF EXISTS github_stats;  
DROP TABLE IF EXISTS twitch_stats;
DROP TABLE IF EXISTS leetcode_stats;