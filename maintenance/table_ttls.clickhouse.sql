alter table system.asynchronous_metric_log
    modify TTL event_time + INTERVAL 1 DAY DELETE;

alter table system.metric_log
    modify TTL event_time + INTERVAL 1 WEEK DELETE,
               event_time + INTERVAL 1 DAY RECOMPRESS CODEC(LZ4HC(10));

alter table system.query_log
    modify TTL event_time + INTERVAL 1 MONTH DELETE,
               event_time + INTERVAL 1 WEEK RECOMPRESS CODEC(LZ4HC(10));

alter table system.trace_log
    modify TTL event_time + INTERVAL 1 WEEK DELETE,
               event_time + INTERVAL 1 DAY RECOMPRESS CODEC(LZ4HC(10));

alter table system.part_log
    modify TTL event_time + INTERVAL 1 MONTH DELETE,
               event_time + INTERVAL 1 WEEK RECOMPRESS CODEC(LZ4HC(10));

