# Web Vitals Collector

A blazing fast web application to capture web vital signals from the javascript library 
[web-vitals](https://github.com/GoogleChrome/web-vitals) and store the in a 
[clickhouse database](https://clickhouse.com/)

## Configuration

Configuration is handled by setting environment variables

* `CLICKHOUSE_HOST`: Hostname of the clickhouse server (default: "localhost")
* `CLICKHOUSE_PORT`: Port of the clickhouse database (default: "9000")
* `CLICKHOUSE_DATABASE`: Name of the clickhouse database (default: "default")
* `CLICKHOUSE_USERNAME`: Username for the clickhouse database (default: "default")
* `CLICKHOUSE_PASSWORD`: Password for the clickhouse database (default: "")
* `PORT`: App port (default: 3000)


## Instrument the clients

To send web vitals to this application, we use the [web-vitals](https://github.com/GoogleChrome/web-vitals) 
javascript library.

### Instrumentation example

```javascript
import { getCLS, getFID, getLCP, getFCP, getTTFB } from 'web-vitals';

const COLLECTOR_ENDPOINT = 'http://localhost:3000/';

const queue = new Set();
function addToQueue(metric) {
    queue.add(metric);
}

function flushQueue() {
    if (queue.size <= 0) {
        return;
    }

    const body = Object.fromEntries(
        [...queue]
            .filter((entry) => entry.name && entry.value)
            .map((entry) => {
                return [entry.name.toLowerCase(), entry.value];
            })
    );

    body.url = window.location.href;

    if (navigator.sendBeacon) {
        navigator.sendBeacon(COLLECTOR_ENDPOINT, JSON.stringify(body));
    } else {
        fetch(COLLECTOR_ENDPOINT, {
            body: JSON.stringify(body),
            method: 'POST',
            keepalive: true,
        });
    }

    queue.clear();
}

getCLS(addToQueue);
getFID(addToQueue);
getLCP(addToQueue);
getFCP(addToQueue);
getTTFB(addToQueue);

window.addEventListener('visibilitychange', () => {
    if (document.visibilityState === 'hidden') {
        flushQueue();
    }
});

// Safari workaround
window.addEventListener('pagehide', flushQueue);
```

## Benchmark

Benchmarks are build with [wrk](https://github.com/wg/wrk), use the following command to run it

```shell
$ wrk -t16 -c400 -d30s -s loadtest/config.lua http://localhost:3000/
```

### Benchmark results

On a 9th generation Intel Core i9, running both this app and the clickhouse database, the following results were 
observed:

```
Running 30s test @ http://localhost:3000/
  16 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    18.51ms   47.05ms 251.50ms   89.09%
    Req/Sec    20.56k     7.85k   75.10k    75.27%
  8946184 requests in 30.09s, 0.98GB read
Requests/sec: 297330.52
Transfer/sec:     33.46MB
```