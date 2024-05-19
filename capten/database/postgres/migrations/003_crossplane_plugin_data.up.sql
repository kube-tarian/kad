CREATE TABLE IF NOT EXISTS crossplane_project (
    id INTEGER PRIMARY KEY,
    git_project_id UUID,
    git_project_url varchar(200),
    status varchar(40),
    last_update_time TIMESTAMP
);

CREATE TABLE IF NOT EXISTS crossplane_provider (
    id UUID PRIMARY KEY,
    cloud_provider_id varchar(80),
    provider_name varchar(40),    
    cloud_type varchar(40),    
    status varchar(40),
    last_update_time TIMESTAMP
);

CREATE TABLE IF NOT EXISTS managed_cluster (
    id UUID PRIMARY KEY,
    cluster_name varchar(80),
    cluster_endpoint varchar(200),
    cluster_deploy_status varchar(40),
    app_deploy_status varchar(40),
    last_update_time TIMESTAMP
);