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

async producer 10000 times
![image](https://github.com/user-attachments/assets/4ba4b772-f2d4-4c0f-bc77-04c9ca1340e6)

sync producer 10000 times
![image](https://github.com/user-attachments/assets/5c3652bb-829b-4183-ac13-9c2971e8cdb6)

