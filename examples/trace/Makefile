.PHONY: proto

grpcsrv:
	@go run grpc-service/main.go

httpsrv:
	@go run http-service/main.go

httpint:
	@go run http-interface/main.go

http2http:
	@curl -v localhost:9000/http-to-http/1

http2grpc:
	@curl -v localhost:9000/http-to-grpc/1

http2grpc2http:
	@curl -v localhost:9000/http-to-grpc-to-http/1

http2grpc2httperr:
	@curl -v localhost:9000/http-to-grpc-to-http/2


