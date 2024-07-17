### required
```bash 
# mysql
mysql> CREATE USER 'exporter'@'%' IDENTIFIED BY '<PASSWORD>' WITH MAX_USER_CONNECTIONS 3;
mysql> GRANT PROCESS, REPLICATION CLIENT, SELECT ON *.* TO 'exporter'@'%';
mysql> flush privileges;
```