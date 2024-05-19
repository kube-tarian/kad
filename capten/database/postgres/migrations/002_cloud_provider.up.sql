CREATE TABLE IF NOT EXISTS cloud_provider (
    id UUID PRIMARY KEY,
    cloud_type varchar(40),
    labels TEXT[],
    last_update_time TIMESTAMP
);