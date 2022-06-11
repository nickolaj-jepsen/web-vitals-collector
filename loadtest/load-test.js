import http from 'k6/http';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

export default function () {
    const payload = JSON.stringify({
        url: `/${randomString(50)}`,
        cls: 10,
        fcp: 100,
        fid: 1000,
        lcp: 10000,
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    http.post('http://localhost:3000', payload, params);
}