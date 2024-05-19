CREATE TABLE IF NOT EXISTS git_project (
    id UUID PRIMARY KEY,
    project_url varchar(200),
    labels TEXT[],
    last_update_time TIMESTAMP
);

CREATE TABLE IF NOT EXISTS cloud_provider (
    id UUID PRIMARY KEY,
    cloud_type varchar(40),
    labels TEXT[],
    last_update_time TIMESTAMP
);

CREATE TABLE IF NOT EXISTS container_registry (
    id UUID PRIMARY KEY,
    registry_url varchar(200),
    registry_type varchar(40),
    labels TEXT[],
    last_update_time TIMESTAMP
);