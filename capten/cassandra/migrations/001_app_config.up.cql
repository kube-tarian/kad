CREATE TABLE ClusterAppConfig(
    app_name TEXT, description TEXT, category TEXT,
    chart_name TEXT, repo_name TEXT, repo_url TEXT,
    namespace TEXT, release_name TEXT, version TEXT,
    launch_url TEXT, launch_redirect_url TEXT,
    create_namespace BOOLEAN, privileged_namespace BOOLEAN,
    override_values TEXT, launch_ui_values TEXT,
    template_values TEXT, default_app BOOLEAN,
    icon BLOB, install_status TEXT,
    update_time TEXT, plugin_name TEXT, plugin_description TEXT,
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
