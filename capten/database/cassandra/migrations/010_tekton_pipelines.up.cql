CREATE TABLE TektonPipelines (
    id UUID,
    pipeline_name TEXT,
    git_org_id TEXT,
    container_registry_id TEXT,
    status TEXT,
    last_update_time TEXT,
    workflow_id TEXT,
    workflow_status TEXT,
    PRIMARY KEY (id)
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