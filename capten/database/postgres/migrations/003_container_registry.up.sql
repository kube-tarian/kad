CREATE TABLE IF NOT EXISTS container_registry (
    id UUID PRIMARY KEY,
    registry_url varchar(200),
    registry_type varchar(40),
    labels TEXT[],
    last_update_time TIMESTAMP
);