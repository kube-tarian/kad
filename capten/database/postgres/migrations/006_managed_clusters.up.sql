CREATE TABLE IF NOT EXISTS managed_clusters (
    id UUID PRIMARY KEY,
    cluster_name varchar(80),
    cluster_endpoint varchar(200),
    cluster_deploy_status varchar(40),
    app_deploy_status varchar(40),
    last_update_time TIMESTAMP
);