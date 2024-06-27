# gounread

no auth
```bash 
$ docker run --name testscylla -d -p 9042:9042 scylladb/scylla --smp 1 --authenticator AllowAllAuthenticator
```
cqlsh 접속
```bash 
$ docker exec -it testscylla cqlsh
```

```bash 
$ docker run --name testscylla -d -p 9042:9042 scylladb/scylla --smp 1 --authenticator PasswordAuthenticator
```
cqlsh 접속(default superuser)
```bash 
$ docker exec -it testscylla cqlsh -u cassandra -p cassandra
```

default superuser는 고정이라 새로운 유저를 만들고 이를 삭제해야한다.
```bash 
$ docker exec -it testscylla cqlsh -u cassandra -p cassandra
cassandra@cqlsh> CREATE ROLE scylla WITH PASSWORD = '1234' AND LOGIN = true AND SUPERUSER = true;
cassandra@cqlsh> LIST ROLES;
 
 role      | super | login | options
-----------+-------+-------+---------
 cassandra |  True |  True |        {}
    scylla |  True |  True |        {}
```

scylla라는 유저에게 superuser를 주었으니 default superuser(cassandra)를 삭제하자
```bash 
$ docker exec -it testscylla cqlsh -u scylla -p 1234
scylla@cqlsh> DROP ROLE cassandra;
scylla@cqlsh> LIST ROLES; 

 role   | super | login | options
--------+-------+-------+---------
 scylla |  True |  True |        {}
```

node status
```bash 
$ docker compose exec -it scylla nodetool status
```

```bash 
Datacenter: datacenter1
=======================
Status=Up/Down
|/ State=Normal/Leaving/Joining/Moving
--  Address     Load       Tokens       Owns    Host ID                               Rack
UN  172.26.0.2  418.35 KB  256          ?       1ca8e6af-9233-405f-b4a2-75c6b603e9ca  rack1
UN  172.26.0.3  400.8 KB   256          ?       a8098db1-f596-4f9a-9701-6a39daff00f6  rack1
UN  172.26.0.4  412.41 KB  256          ?       c97a498f-f354-48c9-a5dd-fdf9855402c2  rack1

Note: Non-system keyspaces don't have the same replication settings, effective ownership information is meaningless
```
thread safe
`The Java List Index is not thread safe. The set or map collection types are safer for updates.`
see https://docs.datastax.com/en/cql-oss/3.3/cql/cql_reference/cqlUpdate.html#Updatingalist

