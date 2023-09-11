mock:
	mockgen -destination service/mock/service.go -package mocksrv github.com/YongJeong-Kim/go/goasynq/service Servicer
	mockgen -destination store/mock/query.go -package mockq github.com/YongJeong-Kim/go/goasynq/store Store
	mockgen -destination worker/mock/distributor.go -package mockwk github.com/YongJeong-Kim/go/goasynq/worker Distributor

test:
	go test ./...