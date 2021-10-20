# gogin

[viper](https://github.com/spf13/viper) 설치

서버포트 8081로 변경하기

```bash
$ make server SERVER_ADDRESS=0.0.0.0:8081
```

[gomock](https://github.com/golang/mock) 설치

```bash
$ mockgen -package <mock package name> -destination <mock go file> <Store 경로> <Store interface name>
$ mockgen -package mockdb -destination db/mock/store.go github.com/yongjeong-kim/go/gogin/db/sqlc Store
```

네트워크 생성
```bash
$ docker network create <network_name>
$ docker network create gin-network
```

mariadb 네트워크 연결하기
```bash
$ docker network connect <container_name>
$ docker network connect mariadb
```

mariadb와 같은 네트워크로 연결하기
```bash
$ docker run --name <container_name> --network <network_name> -p <host_ports:container_ports> -e GIN_MODE=release -e DB_SOURCE="<DB_USER>:<DB_PASSWORD>@tcp(<container_name>:<container_ports>)/go?parseTime=true" <image_name>
$ docker run --name ginserver --network gin-network -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="root:1234@tcp(mariadb:3306)/go?parseTime=true" <image_name>  
```