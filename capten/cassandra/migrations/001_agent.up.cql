CREATE TABLE AppConfig(
    app_name text, description text, category text,
    chart_name text, repo_name text, repo_url text,
    namespace text, release_name text, version text,
    launch_url text, launch_redirect_url text,
    create_namespace boolean, privileged_namespace boolean,
    override_values text, launch_ui_values text,
    icon blob, install_status text,
    PRIMARY KEY (release_name)
) WITH bloom_filter_fp_chance = 0.01
    AND caching = {'keys': 'ALL', 'rows_per_partition': 'NONE'}
    AND comment = ''
    AND compaction = {'class': 'org.apache.cassandra.db.compaction.LeveledCompactionStrategy', 'tombstone_compaction_interval': '1800', 'tombstone_threshold': '0.01', 'unchecked_tombstone_compaction': 'true'}
    AND compression = {'chunk_length_in_kb': '64', 'class': 'org.apache.cassandra.io.compress.LZ4Compressor'}
    AND crc_check_chance = 1.0
    AND default_time_to_live = 0
    AND gc_grace_seconds = 3600
    AND max_index_interval = 2048
    AND memtable_flush_period_in_ms = 0
    AND min_index_interval = 128
    AND speculative_retry = '99PERCENTILE';