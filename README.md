# Redis-Lite

A basic Redis server built from scratch in Go.

## Features
- RESP protocol parser
- kqueue-based async event loop (macOS)
- Supported commands: PING, SET, GET

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

## Roadmap
- EXPIRE, TTL, INCR, DECR
- List, Hash, Set commands
- Persistence (RDB/AOF)
