CREATE TABLE IF NOT EXISTS crossplane_provider (
    id UUID PRIMARY KEY,
    cloud_provider_id varchar(80),
    provider_name varchar(40),    
    cloud_type varchar(40),    
    status varchar(40),
    last_update_time TIMESTAMP
);