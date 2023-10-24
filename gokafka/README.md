#### create topic
```bash
$ docker-compose exec kafka kafka-topics --create --topic <TOPIC_NAME> --bootstrap-server kafka:9092 --replication-factor 1 --partitions 1
```

```bash 
$ docker-compose exec kafka kafka-topics --describe --topic <TOPIC_NAME> --bootstrap-server kafka:9092 
```

#### consumer
```bash
$ docker compose exec kafka bash
$ kafka-console-consumer --topic <TOPIC_NAME> --bootstrap-server kafka:9092
```

#### producer
```bash 
$ docker compose kafka bash 
$ kafka-console-producer --topic <TOPIC_NAME> --broker-list kafka:9092
```