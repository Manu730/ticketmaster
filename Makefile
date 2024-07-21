# Makefile

protoc:
	cd protos && protoc --go_out=../protogen/golang --go_opt=paths=source_relative --go-grpc_out=../protogen/golang --go-grpc_opt=paths=source_relative  ticketmaster.proto
test:
	go test ./... -v

run_server:
	cd server && go run server.go

run_client:
	cd client && go run client.go

build_server:
	cd server && GOOS=linux go build -o ticketmasterserver server.go

build_client:
	cd client && GOOS=linux go build -o ticketmasterclient client.go

build_server_docker:
	cd server && docker build -t ticketmasterserver:v1 .

build_client_docker:
	cd client && docker build -t ticketmasterclient:v1 .
