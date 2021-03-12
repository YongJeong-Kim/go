# go sqlc

### [golang-sqlc](https://github.com/kyleconroy/sqlc) 설치하기

현재 sqlc 1.7.0 버전 기준 window에서 postgresql이 동작하지 않아 mysql로 진행.

### docker mysql 설치
```bash
$ docker pull mysql
```

``` bash
no matching manifest for windows/amd64 10.0.18363 in the manifest list entries
```
라는 에러가 뜬다면

docker - Settings - Docker Engine 에서 아래의 내용을
```json
{
  "registry-mirrors": [],
  "insecure-registries": [],
  "debug": false,
  "experimental": false
}
```
```json
{
  ...,
  "experimental": true
}
```
로 수정한 뒤 restart

### docker mysql image 생성
```bash
$ docker run --name <some_mysql_name> -e MYSQL_ROOT_PASSWORD=<password> -p <host_ports:container_ports> -d mysql:<tag>
$ docker run --name mysql8 -e MYSQL_ROOT_PASSWORD=1234 -p 3306:3306 -d mysql:latest
```

### docker mysql container 접속
```bash
$ docker exec -it <container_name> bash
```

```bash
$ docker exec -it mysql8 bash
$ mysql -u root -p
// or
$ docker exec -it mysql8 mysql -u root -p
```

#### database 생성
```bash
$ create database <database_name> default character set utf8 collate utf8_general_ci
```

### sqlc.yaml 설정 파일 생성하기
```bash
$ sqlc init
```