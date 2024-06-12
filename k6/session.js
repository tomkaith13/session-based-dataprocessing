
import { Trend } from 'k6/metrics';
import http from 'k6/http'
import sleep from 'k6'
import { check } from 'k6';

export let options = { maxRedirects: 4, 
  // duration: '20s' 
    
  scenarios: {
    // constant_request_rate_100rps: {
    //   executor: 'constant-arrival-rate',
    //   rate: 100,
    //   timeUnit: '10s',
    //   duration: '3m',
    //   preAllocatedVUs: 5,
    //   // maxVUs: 50,
    //   env: { EXAMPLEVAR: 'create_thread_v2' },
    // },
    // constant_request_rate_50rps: {
    //   executor: 'constant-arrival-rate',
    //   rate: 50,
    //   timeUnit: '1s',
    //   duration: '1m',
    //   preAllocatedVUs: 10,
    //   env: { EXAMPLEVAR: 'create_thread_v2' },
    // },
    constant_request_rate_10rps: {
      executor: 'constant-arrival-rate',
      rate: 10,
      timeUnit: '1s',
      duration: '1m',
      preAllocatedVUs: 10,
      env: { EXAMPLEVAR: 'create_thread_v2' },
    },
  },
};


const trend1 = new Trend('filter using mongo', true);
const trend2 = new Trend('filter using parquet', true);
const trend3 = new Trend('filter using parquet from table', true);


let customMetrics = {};
for (let key in options.scenarios) {
  options.scenarios[key].env['MY_SCENARIO'] = key;
  let customMetricName = key + '_' + options.scenarios[key].env.EXAMPLEVAR;
  customMetrics[key] = new Trend(customMetricName, true);
}

export function setup() {
  const setup_resp1 = http.post('http://localhost:8080/person');
  check(setup_resp1, {
    'setup person is status 200': (r) => r.status === 200,
  });

  const setup_resp2 = http.post('http://localhost:8080/person-parquet');
  check(setup_resp2, {
    'setup person-parquet is status 201': (r) => r.status === 201,
  });

  const setup_resp3 = http.post('http://localhost:8080/person-parquet-table');
  check(setup_resp3, {
    'setup person-parquet-table is status 201': (r) => r.status === 202,
  });
}

export default function () {
  const response = http.post('http://localhost:8080/person-search');

  check(response, {
    'person-search is status 200': (r) => r.status === 200,
  });

  if (response.status == 200) {
    trend1.add(response.timings.duration);
    // customMetrics[__ENV.MY_SCENARIO].add(response.timings.duration);
  }


  const response2 = http.post('http://localhost:8080/person-search-parquet');

  check(response2, {
    'person-parquet-search is status 200': (r) => r.status === 200,
  });

  if (response2.status == 200) {
    trend2.add(response2.timings.duration);
    // customMetrics[__ENV.MY_SCENARIO].add(response2.timings.duration);
  }


  const response3 = http.post('http://localhost:8080/person-search-parquet-table');

  check(response3, {
    'person-parquet-search-table is status 200': (r) => r.status === 200,
  });

  if (response3.status == 200) {
    trend3.add(response3.timings.duration);
    // customMetrics[__ENV.MY_SCENARIO].add(response3.timings.duration);
  }
}