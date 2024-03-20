mockgen:
	mockgen -package mock -destination service/mock/service.go github.com/yongjeong-kim/go/goapigw/account/service Servicer

mockgenservice:
	mockgen -package mock -destination service/mock/service.go -source service/service.go Servicer

mockgentoken:
	mockgen -packge tkmock -destination token/mock/token.go -source token/token.go TokenMaker

test:
	go test ./...