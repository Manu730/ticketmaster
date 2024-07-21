FROM ubuntu:latest
LABEL Description="Grpc server for ticket booking"
ADD ticketmasterserver /ticketmasterserver
ENTRYPOINT ["./ticketmasterserver"]
