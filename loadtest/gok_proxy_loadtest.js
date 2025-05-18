import http from "k6/http";
import { check, sleep, group } from "k6";
import { Trend } from "k6/metrics";
import { htmlReport } from "https://raw.githubusercontent.com/benc-uk/k6-reporter/main/dist/bundle.js";
import { textSummary } from "https://jslib.k6.io/k6-summary/0.0.1/index.js";

// --- Configuration ---
const PROXY_URL = "http://localhost:8080"; // Your gok-proxy address
const TARGET_HTTP_URL = "http://httpbin.org/get"; // Target for HTTP GET requests
const TARGET_HTTPS_URL_FOR_CONNECT_TEST = "https://httpbin.org/get"; // Target for CONNECT test (will be proxied)

// --- Custom Metrics ---
const httpGetDuration = new Trend("http_get_duration", true);
const httpConnectAndGetDuration = new Trend(
  "http_connect_and_get_duration",
  true
);

// --- k6 Options ---
export const options = {
  // Configure k6 to use gok-proxy for all HTTP requests
  // This is one way; k6 also respects HTTP_PROXY/HTTPS_PROXY env vars.
  // However, for explicit control within the script, we'll pass the proxy URL to http.get/request.
  // Note: k6's built-in HTTP client handles proxying for standard requests.
  // For CONNECT, we manually craft the request.

  stages: [
    { duration: "10s", target: 50 }, // Ramp up to 20 virtual users over 30s
    { duration: "15s", target: 100 }, // Stay at 20 virtual users for 1m
    { duration: "5s", target: 0 }, // Ramp down to 0 users
  ],
  thresholds: {
    http_req_failed: ["rate<0.01"], // http errors should be less than 1%
    http_req_duration: ["p(95)<1500"], // 95% of requests should be below 1500ms (increased for proxied https)
    http_get_duration: ["p(95)<800"],
    http_connect_and_get_duration: ["p(95)<1200"], // Duration for the whole HTTPS (CONNECT + GET) transaction
  },
};

// --- Setup Function (runs once before VU init) ---
export function setup() {
  console.log(
    "k6 load test starting. Ensure gok-proxy is running and accessible."
  );
  console.log(`Targeting gok-proxy assumed to be at: ${PROXY_URL}`);
  console.log(
    "To test through the proxy, ensure you run k6 with the HTTP_PROXY environment variable set:"
  );
  console.log(
    `  HTTP_PROXY=${PROXY_URL} k6 run loadtest/gok_proxy_loadtest.js`
  );
  console.log("---");
  console.log(
    `HTTP GET requests will be sent to: ${TARGET_HTTP_URL} (via proxy)`
  );
  console.log(
    `HTTPS (CONNECT then GET) requests will be sent to: ${TARGET_HTTPS_URL_FOR_CONNECT_TEST} (via proxy)`
  );
  console.log("---");
}

// --- Main Test Function (Default Scenario for each VU) ---
export default function () {
  // Group for HTTP GET requests
  group("http_get_proxied_requests", function () {
    // k6 uses HTTP_PROXY env var. No special params needed here for http.get.
    const res = http.get(TARGET_HTTP_URL);

    check(res, {
      "[HTTP GET] status is 200": (r) => r.status === 200,
      "[HTTP GET] body is not empty": (r) => r.body && r.body.length > 0,
    });
    httpGetDuration.add(res.timings.duration);
  });

  sleep(0.5); // Simulate some think time

  // Group for HTTPS (CONNECT then GET) requests
  group("https_connect_then_get_proxied_requests", function () {
    // When HTTP_PROXY is set and an HTTPS URL is requested, k6 automatically
    // sends a CONNECT request to the proxy, then establishes TLS to the target,
    // then sends the GET request through the tunnel.
    const res = http.get(TARGET_HTTPS_URL_FOR_CONNECT_TEST);

    check(res, {
      "[HTTPS via CONNECT] status is 200": (r) => r.status === 200,
      "[HTTPS via CONNECT] body is not empty": (r) =>
        r.body && r.body.length > 0,
    });
    // This duration includes: CONNECT to proxy, proxy connects to target,
    // TLS handshake through proxy, HTTP GET through proxy, response.
    httpConnectAndGetDuration.add(res.timings.duration);
  });

  sleep(1); // Simulate more think time
}

// --- Teardown Function (runs once after all VUs complete) ---
export function teardown(data) {
  console.log("k6 load test finished.");
}

// --- Handle Summary (Generate report at the end) ---
export function handleSummary(data) {
  console.log("Preparing k6 summary...");
  return {
    stdout: textSummary(data, { indent: " ", enableColors: true }), // ✔ fixed here
    "loadtest/summary.html": htmlReport(data),
  };
}
