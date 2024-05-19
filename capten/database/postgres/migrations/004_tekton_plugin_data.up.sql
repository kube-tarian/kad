CREATE TABLE IF NOT EXISTS tekton_project (
    id INTEGER PRIMARY KEY,
    git_project_id UUID,
    git_project_url varchar(200),
    status varchar(40),
    last_update_time TIMESTAMP
);