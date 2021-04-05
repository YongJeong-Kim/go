#gogin

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