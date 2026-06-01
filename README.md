# Redis-Lite

A basic Redis server built from scratch in Go.

## Features
- RESP protocol parser
- kqueue-based async event loop (macOS)
- Supported commands for now: PING, SET, GET, TTL, DEL, EXPIRE
- Supports ACTIVE and PASSIVE mode key deletion

## Run

```bash
go run main.go
```

## Test

```bash
redis-cli -p 7379 PING
redis-cli -p 7379 SET foo bar
redis-cli -p 7379 GET foo
```