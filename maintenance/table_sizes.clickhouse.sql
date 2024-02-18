SELECT table, formatReadableSize(size) as size, rows, days, formatReadableSize(avgDaySize) as avgDaySize
FROM (
    SELECT
    table, sum (bytes) AS size, sum (rows) AS rows, min (min_date) AS min_date, max (max_date) AS max_date, (max_date - min_date) AS days, size / (max_date - min_date) AS avgDaySize
    FROM system.parts
    WHERE active
    GROUP BY table
    ORDER BY rows DESC
    );

SELECT name, table, removal_state, bytes, last_removal_attemp_time, min_date, max_date FROM system.parts
WHERE min_date < '2024-01-01'
LIMIT 30;

