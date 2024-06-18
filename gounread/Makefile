
testall:
	ginkgo -r

mocksvc:
	mockgen -source service/service.go -destination service/svcmock/service_mock.go -package svcmock Servicer