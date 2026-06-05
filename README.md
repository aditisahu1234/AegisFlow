# AegisFlow

## Adaptive Distributed API Protection Platform

AegisFlow is a backend infrastructure project focused on protecting APIs from abuse, traffic spikes, and malicious requests.

The long-term goal is to build an intelligent API protection layer inspired by modern systems such as Cloudflare and enterprise API gateways, featuring distributed rate limiting, observability, adaptive traffic control, and AI-driven traffic intelligence.

This repository currently contains **Phase 1: API Protection Layer**, implemented in Go.

---

## Current Features

### Sliding Window Rate Limiter

* Per-IP request tracking
* Sliding Window algorithm
* Configurable request limits
* HTTP 429 (Too Many Requests) response when limits are exceeded

### Middleware-Based Protection

All incoming requests pass through a protection layer before reaching application endpoints.

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

### Metrics Endpoint

```http
GET /metrics
```

Returns:

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

### Load Testing

Validated using K6 load testing.

---

## Architecture

```text
Client
   ↓
Sliding Window Middleware
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

Start the server:

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

## Roadmap

### Level 2

* Redis-backed Distributed Rate Limiting
* Shared State Across Multiple Nodes
* Multi-instance Support

### Level 3

* Prometheus Metrics
* Grafana Dashboards
* Failure Recovery
* Dynamic Throttling

### Level 4

* Traffic Intelligence Engine
* Threat Classification
* Adaptive Rate Limits
* AI-driven Decisions

### Level 4.5

* Docker
* Kubernetes (K3s)
* Horizontal Pod Autoscaling
* OpenTelemetry Tracing
* GitHub Actions CI/CD

---

## Tech Stack

* Go
* net/http
* K6
* Git
