# gonats

[nats](https://docs.nats.io/)

[default port](https://docs.nats.io/running-a-nats-service/nats_docker#usage)
- 4222 is for clients.
- 8222 is an HTTP management port for information reporting.
- 6222 is a routing port for clustering.

```bash
$ docker compose up
```

#### monitoring
- localhost:8222


#### subscribe
```bash 
$ go run cmd/sub/subscribe.go
```

#### publish
```bash 
$ go run cmd/pub/publish.go
```

max connections default 65536

see [monitoring](https://docs.nats.io/running-a-nats-service/nats_admin/monitoring#arguments-1)
`localhost:8222/connz`
default
```json
{ 
  "server_id": "...",
  "now": "...",
  "num_connections": 0,
  "total": 0,
  "offset": 0,
  "limit": 1024, // Number of results to return. Default is 1024.
  "connections": [...]
}
```

1028 connections json
```json
{
  "server_id": "...",
  "now": "...",
  "num_connections": 1024,
  "total": 1028,
  "offset": 0,
  "limit": 1024, // Number of results to return. Default is 1024.
  "connections": [...]
}
```

localhost:8222/connz?limit=1028
```json 
{
  "server_id": "...",
  "now": "...",
  "num_connections": 1028,
  "total": 1028,
  "offset": 0,
  "limit": 1028, // Number of results to return. Default is 1024.
  "connections": [...]
}
```
