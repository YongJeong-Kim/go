# go1m

create topic
```bash 
$ docker compose exec kafka kafka-topics.sh --create --topic aaa --bootstrap-server localhost:19092 --replication-factor 3 --partitions 3
```

describe topic
```bash 
$ docker compose exec kafka kafka-topics.sh --describe --topic aaa --bootstrap-server localhost:19092
```

produce 
```bash 
$ docker compose exec kafka kafka-console-producer.sh --topic bbb --bootstrap-server localhost:19092
```

produce round robin
```bash 
$ docker compose exec kafka kafka-console-producer.sh --topic bbb --bootstrap-server localhost:19092 --producer-property partitioner.class=org.apache.kafka.clients.producer.RoundRobinPartitioner
```

consume 
```bash
$ docker compose exec kafka kafka-console-consumer.sh --topic aaa --bootstrap-server localhost:19092 --partition 0
```

consume group
```bash
$ docker compose exec kafka kafka-console-consumer.sh --topic aaa --bootstrap-server localhost:19092 --group group111
```

consumer group list
```bash
$ docker compose exec kafka kafka-consumer-groups.sh --bootstrap-server localhost:19092 --list
```

describe group
```bash
$ docker compose exec kafka kafka-consumer-groups.sh --bootstrap-server localhost:19092 --describe --group group11
```