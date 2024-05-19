CREATE TABLE IF NOT EXISTS git_project (
    id UUID PRIMARY KEY,
    project_url varchar(200),
    labels TEXT[],
    last_update_time TIMESTAMP
);