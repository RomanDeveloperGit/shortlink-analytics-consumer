-- +goose Up
CREATE TABLE link_visited_event
(
    event_id        UUID,

    link_id         Int,
    short_code      String,
    full_url        String,
    timestamp       DateTime,

    -- ip              IPv6,
    -- user_agent      String,
    -- referer         String DEFAULT '',
    -- country         LowCardinality(String),
    -- city            String DEFAULT '',

    -- Партиционирование
    event_date      Date DEFAULT toDate(timestamp)
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(event_date)
ORDER BY (event_date, link_id)
TTL event_date + INTERVAL 2 YEAR
SETTINGS index_granularity = 8192;

-- +goose Down
DROP TABLE IF EXISTS link_visited_event;
