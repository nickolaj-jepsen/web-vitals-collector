# Web Vitals Collector

A blazing fast web application to capture web vital signals from the javascript library [web-vitals](https://github.com/GoogleChrome/web-vitals) and store the in a [clickhouse database](https://clickhouse.com/)

## Configuration

Configuration is handled by enviroment variables

* `CLICKHOUSE_HOST`: Hostname of the clickhouse database (default: "localhost")
* `CLICKHOUSE_PORT`: Port of the clickhouse database (default: "9000")
* `CLICKHOUSE_DATABASE`: Database name (default: "default")
* `CLICKHOUSE_USERNAME`: Username used to connect to the clickhouse server (default: "default")
* `CLICKHOUSE_PASSWORD`: Password used to connect to the clickhouse server (default: "")
* `PORT`: App port (default: 3000)


## Benchmark

Benchmarks are build with [k6](https://k6.io), use the following command to run it

```shell
$ k6 run --vus 10 --duration 30s loadtest/load-test.js
```
