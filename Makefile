# Makefile

protoc:
	cd protos && protoc --go_out=../protogen/golang --go_opt=paths=source_relative --go-grpc_out=../protogen/golang --go-grpc_opt=paths=source_relative  ticketmaster.proto
test:
	go test ./... -v
