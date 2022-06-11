CREATE TABLE vitals (
    `timestamp`  DateTime,
    `url`        String,
    `identifier` String,
    `cls`        Nullable(Float64),
    `fcp`        Nullable(Float64),
    `fid`        Nullable(Float64),
    `lcp`        Nullable(Float64),
    `ttfb`       Nullable(Float64)
) ENGINE = MergeTree
ORDER BY (toStartOfDay(`timestamp`), path(`url`))
    TTL `timestamp` + INTERVAL 7 DAY GROUP BY toStartOfDay(`timestamp`), path(`url`) SET `cls` = avg(`cls`),
                                                                                         `fcp` = avg(`fcp`),
                                                                                         `fid` = avg(`fid`),
                                                                                         `lcp` = avg(`lcp`),
                                                                                         `ttfb` = avg(`ttfb`),
                                                                                         `identifier` = any(`identifier`)