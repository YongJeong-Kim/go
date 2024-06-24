
testall:
	ginkgo -r

mocksvc:
	mockgen -source service/service.go -destination service/svcmock/service_mock.go -package svcmock Servicer

join1u1:
	go run client/room/recv/receive.go --user 05f84f46-d4ad-42db-af4f-da63cffcb721 --room 01f84cfa-e487-494c-82e5-e75f95ef0573

join1u2:
	go run client/room/recv/receive.go --user febba554-152e-496a-add5-31d0672fdc2a --room 01f84cfa-e487-494c-82e5-e75f95ef0573

join1u3:
	go run client/room/recv/receive.go --user 01a5482c-2053-4728-96b0-8c950ebcf35a --room 01f84cfa-e487-494c-82e5-e75f95ef0573