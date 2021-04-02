#gogin

[gomock](https://github.com/golang/mock) 설치


```bash
$ mockgen -package <mock package name> -destination <mock go file> <Store 경로> <Store interface name>
$ mockgen -package mockdb -destination db/mock/store.go github.com/yongjeong-kim/go/gogin/db/sqlc Store
```