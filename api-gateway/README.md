# AegisFlow

## Adaptive Distributed API Protection Platform

AegisFlow is a production-oriented backend infrastructure platform designed to protect APIs from abuse, traffic spikes, dependency failures, and infrastructure instability.

Inspired by systems such as Cloudflare, modern API gateways, and reliability engineering practices used across large-scale distributed systems, AegisFlow combines distributed rate limiting, resilience engineering, observability, adaptive traffic control, and AI-driven traffic intelligence.

The long-term vision is to build an intelligent protection platform capable of automatically detecting threats, adapting traffic policies, surviving infrastructure failures, and making operational decisions based on real-time system conditions.

---

## Current Features

### Distributed Sliding Window Rate Limiter

Built from scratch using Redis Sorted Sets and Atomic Lua Scripts.

Features:

* Per-IP request tracking
* Accurate Sliding Window algorithm
* Millisecond precision
* Configurable request limits
* HTTP 429 (Too Many Requests) enforcement
* Shared distributed state through Redis
* Atomic rate-limiting operations

---

### Dependency Timeout Protection

Every Redis operation executes with bounded latency using Go contexts.

Features:

* Fast failure detection
* Prevents hanging requests
* Limits dependency latency impact
* Protects application responsiveness

Implementation:

```go
context.WithTimeout(...)
```

---

### Retry + Exponential Backoff + Full Jitter

Inspired by Amazon reliability engineering patterns.

Features:

* Controlled retry behaviour
* Exponential backoff
* AWS-style Full Jitter
* Retry storm mitigation
* Reduced dependency pressure during outages

Backoff Formula:

```text
sleep = random(
    0,
    min(cap, base * 2^attempt)
)
```

---

### Redis Health Monitoring

Continuous background monitoring of Redis availability.

Features:

* Automatic outage detection
* Automatic recovery detection
* Shared health state
* Infrastructure visibility

---

### Production Circuit Breaker

Implemented using:

```text
github.com/sony/gobreaker/v2
```

Features:

* Closed → Open → Half-Open → Closed state machine
* Rolling failure tracking
* Failure-rate-based tripping
* Automatic recovery testing
* Fail-fast protection
* Dependency isolation

Current Configuration:

* 10-second rolling window
* 1-second buckets
* Failure-rate threshold detection
* Half-open recovery validation
* Automatic state transitions

---

### Middleware-Based Protection

All incoming requests pass through the protection layer before reaching application endpoints.

Protection Pipeline:

```text
Request
    ↓
Rate Limiter
    ↓
Resilience Layer
    ↓
Application Endpoint
```

---

### Health Endpoint

```http
GET /health
```

Returns service health status.

Example:

```json
{
  "status": "healthy"
}
```

---

### Metrics Endpoint

```http
GET /metrics
```

Currently exposes:

* Total Requests
* Allowed Requests
* Blocked Requests

Example:

```json
{
  "total_requests": 150,
  "allowed_requests": 100,
  "blocked_requests": 50
}
```

Future versions will expose:

* Circuit Breaker Metrics
* Failure Counts
* Retry Counts
* Timeout Events
* Recovery Events
* Redis Latency Metrics

---

### Load Testing

Validated using K6 load testing.

Benchmarked for:

* Throughput
* Request Rate
* Latency
* Failure Behaviour
* Recovery Behaviour

---

## Architecture

```text
Client

↓

Sliding Window Middleware

↓

Resilience Layer

├── Dependency Timeout Protection
├── Retry + Exponential Backoff
├── AWS Full Jitter
├── Redis Health Monitoring
├── Sony Gobreaker Circuit Breaker

↓

Redis Distributed Rate Limiter

↓

Protected API Endpoint

↓

Response
```

---

## API Endpoints

### Protected Endpoint

```http
GET /api/data
```

### Health Check

```http
GET /health
```

### Metrics

```http
GET /metrics
```

---

## Running Locally

Start Redis:

```bash
redis-server
```

Start the application:

```bash
go run .
```

Access:

```text
http://localhost:8080/api/data
http://localhost:8080/health
http://localhost:8080/metrics
```

---

## Load Testing

Run:

```bash
k6 run load-test.js
```

---

## Tech Stack

### Backend

* Go
* net/http

### Storage

* Redis
* Redis Sorted Sets
* Lua Scripts

### Reliability Engineering

* sony/gobreaker/v2
* cenkalti/backoff/v4
* context.WithTimeout

### Testing

* K6

### Version Control

* Git
* GitHub

---

## Current Project Status

### Level 2 — Distributed Rate Limiting

✅ Redis-backed Sliding Window

✅ Atomic Lua Scripts

✅ Shared Distributed State

✅ Multi-instance Architecture Foundation

---

### Level 3 — Production Resilience Platform

✅ Dependency Timeouts

✅ Retry Mechanism

✅ Exponential Backoff

✅ AWS Full Jitter

✅ Redis Health Monitoring

✅ Sony Gobreaker Circuit Breaker

✅ Fail-Fast Protection

🔄 Fallback Local Limiter

🔄 Distributed Failure Tracking

🔄 Prometheus Metrics

🔄 Grafana Dashboards

🔄 Adaptive Load Shedding

---

## Roadmap

### Level 4 — Adaptive Intelligent Protection Platform

* Traffic Intelligence Engine
* Threat Classification
* Anomaly Detection
* Adaptive Rate Limits
* Traffic Forecasting
* AI-driven Decisions

---

### Level 4.5 — Platform Engineering Layer

* Docker
* Kubernetes (K3s)
* Horizontal Pod Autoscaling
* OpenTelemetry Tracing
* GitHub Actions CI/CD

---

### Level 5 — Enterprise Reliability

* Chaos Testing
* Multi-Node Failure Simulation
* Distributed Failure Tracking
* Advanced Resilience Policies
* Fault Injection Testing

---

## Project Goal

AegisFlow aims to evolve from a distributed rate limiter into a complete adaptive API protection platform capable of:

* Detecting abnormal traffic
* Protecting backend systems
* Surviving dependency failures
* Automatically adapting protection policies
* Providing full observability
* Leveraging machine learning for operational decision-making

The project is intentionally designed to demonstrate backend engineering, distributed systems, reliability engineering, observability, platform engineering, and applied machine learning in a single production-oriented system.
