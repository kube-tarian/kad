CREATE TABLE IF NOT EXISTS plugin_store_config (
    store_type INTEGER PRIMARY KEY,
    git_project_id varchar(200),
    git_project_url varchar(200),
    status varchar(40),
    last_update_time TIMESTAMP
);

CREATE TABLE IF NOT EXISTS plugin_store_data (
    store_type INTEGER,
    git_project_id varchar(200),
    plugin_name varchar(80),
    category varchar(40),
    versions TEXT[],
    icon BYTEA,
    description TEXT,
    last_update_time TIMESTAMP,
    PRIMARY KEY (store_type, git_project_id, plugin_name)
);
