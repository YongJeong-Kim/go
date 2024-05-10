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