-- +goose Up
CREATE TABLE link_created_event
(
    event_id        UUID,

    link_id         Int,
    short_code      String,
    full_url        String,
    expires_at      DateTime,
    created_at      DateTime,

    -- -- Контекст создания
    -- ip              IPv6,
    -- user_agent      String,
    -- referer         String DEFAULT '',
    -- country         LowCardinality(String),
    -- city            String DEFAULT '',

    -- Партиционирование
    event_date      Date DEFAULT toDate(created_at)
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(event_date)
ORDER BY (event_date, link_id)
TTL event_date + INTERVAL 2 YEAR
SETTINGS index_granularity = 8192;

-- +goose Down
DROP TABLE IF EXISTS link_created_event;
