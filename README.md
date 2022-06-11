# Web Vitals Collector

A blazing fast web application to capture web vital signals from the javascript library [web-vitals](https://github.
com/GoogleChrome/web-vitals) and store the in a [clickhouse database](https://clickhouse.com/)

## Configuration

Configuration is handled by enviroment variables

* `CLICKHOUSE_HOST`: Hostname of the clickhouse database (default: "localhost")
* `CLICKHOUSE_PORT`: Port of the clickhouse database (default: "9000")
* `CLICKHOUSE_DATABASE`: Database name (default: "default")
* `CLICKHOUSE_USERNAME`: Username used to connect to the clickhouse server (default: "default")
* `CLICKHOUSE_PASSWORD`: Password used to connect to the clickhouse server (default: "")
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

Benchmarks are build with [k6](https://k6.io), use the following command to run it

```shell
$ k6 run --vus 10 --duration 30s loadtest/load-test.js
```
