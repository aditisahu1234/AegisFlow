import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
  scenarios: {
    load_test: {
      executor: "constant-vus",
      vus: 100,                  //100 virtual users
      duration: "2m",        //30 seconds of sustained load, instead of fixed number of requests
    },
  },

  //test automatically fails if error rate or latency becomes unacceptable
  thresholds: {
    http_req_failed: ["rate<0.01"],
    http_req_duration: ["p(95)<200"],
  },
};

export default function () {
  const res = http.get("http://localhost:8080/api/data");

  check(res, {
    "status is 200 or 429": (r) =>
      r.status === 200 || r.status === 429,
  });

  sleep(0.1);
}