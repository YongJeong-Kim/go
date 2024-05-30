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

automatically reconnect

