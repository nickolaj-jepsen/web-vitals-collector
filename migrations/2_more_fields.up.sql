ALTER TABLE vitals ADD COLUMN page_type Nullable(String);
ALTER TABLE vitals ADD COLUMN first_load Nullable(Bool);
ALTER TABLE vitals MODIFY TTL `timestamp` + INTERVAL 7 DAY GROUP BY toStartOfDay(`timestamp`), path(`url`) SET `cls` = avg(`cls`),
    `fcp` = avg(`fcp`),
    `fid` = avg(`fid`),
    `lcp` = avg(`lcp`),
    `ttfb` = avg(`ttfb`),
    `page_type` = any(`page_type`),
    `first_load` = min(`first_load`),
    `identifier` = any(`identifier`)