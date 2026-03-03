CREATE TABLE outbox_events (
    id SERIAL PRIMARY KEY,
    event_type VARCHAR(100),
    aggregate_id VARCHAR(100),
    payload JSON,
    status VARCHAR(20),
    retry_count INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    next_retry_at TIMESTAMP NULL
);

CREATE INDEX idx_outbox_events_status_retry_at ON outbox_events(status, next_retry_at);