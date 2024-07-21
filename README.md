# ticketmaster
A grpc server client to book a train ticket

Folder Structure:

        protos ==> Folder for all protobuf definitions
        protogen ==> Folder for generated code from protos, Has sub folders for different language
        protogen/golang ==> generated code from protos for golang
        server   ==>  Grpc server code for ticketmaster app
        client   ==> Grpc client code for ticketmaster app
        Makefile ==> make commands to test/run/build/generate docker images for Grpc server and client

This App needs to be run under GOPATH/src (App assumes golang is already present)

golang code from protobuf definition is already generated

Makefile commands

```

To run UT : make test

To run server: make run_server

To run client: make run_client (to be run on different terminal than server)

To build server: make build_server (builds binary for linux)

To build client: make build_client (builds binary for linux)

To build docker image for server: make build_server_docker

To build docker image for client: make build_client_docker

To generate protobuf code from protos: make protoc

```

Sample Run

UT:

```
ticketmaster % make test                                                            
go test ./... -v
?   	ticketmaster/client	[no test files]
?   	ticketmaster/protogen/golang	[no test files]
=== RUN   TestTicketMaster_BookTicket
=== RUN   TestTicketMaster_BookTicket/user_empty
=== RUN   TestTicketMaster_BookTicket/cost_empty
=== RUN   TestTicketMaster_BookTicket/success


 Received request for ticekt booking for user:  manunero1730@gmail.com
=== RUN   TestTicketMaster_BookTicket/source_station_empty
=== RUN   TestTicketMaster_BookTicket/destination_station_empty
--- PASS: TestTicketMaster_BookTicket (0.01s)
    --- PASS: TestTicketMaster_BookTicket/user_empty (0.00s)
    --- PASS: TestTicketMaster_BookTicket/cost_empty (0.00s)
    --- PASS: TestTicketMaster_BookTicket/success (0.00s)
    --- PASS: TestTicketMaster_BookTicket/source_station_empty (0.00s)
    --- PASS: TestTicketMaster_BookTicket/destination_station_empty (0.00s)
=== RUN   TestTicketMaster_GetReceipt


 Received request for ticekt booking for user:  manunero1730@gmail.com
=== RUN   TestTicketMaster_GetReceipt/success


 Received request for receipt for user:  manunero1730@gmail.com
--- PASS: TestTicketMaster_GetReceipt (0.00s)
    --- PASS: TestTicketMaster_GetReceipt/success (0.00s)
=== RUN   TestTicketMaster_ShowAllocations
=== RUN   TestTicketMaster_ShowAllocations/success


 Received request to show ticket allocations
--- PASS: TestTicketMaster_ShowAllocations (0.07s)
    --- PASS: TestTicketMaster_ShowAllocations/success (0.07s)
=== RUN   TestTicketMaster_RemoveUser
--- PASS: TestTicketMaster_RemoveUser (0.00s)
=== RUN   TestTicketMaster_ModifyUserAllocation
--- PASS: TestTicketMaster_ModifyUserAllocation (0.00s)
PASS
ok  	ticketmaster/server	0.949s
ticketmaster %

```

**Server Run
**

```

ticketmaster % make run_server
cd server && go run server.go
listening on port 8080


 Received request for ticekt booking for user:  manunero1730@gmail.com


 Received request for receipt for user:  manunero1730@gmail.com


 Received request to show ticket allocations


 Received request to show ticket allocations


 Received request to remove user:  manunero1730@gmail.com


 Received request for ticekt booking for user:  manunero1730@gmail.com


 Received request to modify allocation for user:  manunero1730@gmail.com


 Received request to show ticket allocations


 Received request to show ticket allocations

```

**corresponding client run
**

```
ticketmaster % make run_client
cd client && go run client.go


 Recevied booking response:  train_number:1  trip_detail:{source:"London"  destination:"France"}  traveller:{first_name:"Manoj"  last_name:"Kumar"  email:"manunero1730@gmail.com"}


 Recevied receipt:  train_number:1  from_station:"London"  to_station:"France"  user:{first_name:"Manoj"  last_name:"Kumar"  email:"manunero1730@gmail.com"}  cost:{price:20  currency:"Dollar"}


 Recevied allocation data for section A:  train_number:1  section_details:{seats:{seat_no:1  is_allocated:true  allocated_to:{first_name:"Manoj"  last_name:"Kumar"  email:"manunero1730@gmail.com"}}  seats:{seat_no:2}  seats:{seat_no:3}  seats:{seat_no:4}  seats:{seat_no:5}  seats:{seat_no:6}  seats:{seat_no:7}  seats:{seat_no:8}  seats:{seat_no:9}  seats:{seat_no:10}  seats:{seat_no:11}  seats:{seat_no:12}  seats:{seat_no:13}  seats:{seat_no:14}  seats:{seat_no:15}  seats:{seat_no:16}  seats:{seat_no:17}  seats:{seat_no:18}  seats:{seat_no:19}  seats:{seat_no:20}  seats:{seat_no:21}  seats:{seat_no:22}  seats:{seat_no:23}  seats:{seat_no:24}  seats:{seat_no:25}  seats:{seat_no:26}  seats:{seat_no:27}  seats:{seat_no:28}  seats:{seat_no:29}  seats:{seat_no:30}  seats:{seat_no:31}  seats:{seat_no:32}  seats:{seat_no:33}  seats:{seat_no:34}  seats:{seat_no:35}  seats:{seat_no:36}  seats:{seat_no:37}  seats:{seat_no:38}  seats:{seat_no:39}  seats:{seat_no:40}  seats:{seat_no:41}  seats:{seat_no:42}  seats:{seat_no:43}  seats:{seat_no:44}  seats:{seat_no:45}  seats:{seat_no:46}  seats:{seat_no:47}  seats:{seat_no:48}  seats:{seat_no:49}  seats:{seat_no:50}}


 Recevied allocation data for section B:  train_number:1  section_details:{name:B  seats:{seat_no:1  is_allocated:true  allocated_to:{first_name:"Manoj"  last_name:"Kumar"  email:"manunero1730@gmail.com"}}  seats:{seat_no:2}  seats:{seat_no:3}  seats:{seat_no:4}  seats:{seat_no:5}  seats:{seat_no:6}  seats:{seat_no:7}  seats:{seat_no:8}  seats:{seat_no:9}  seats:{seat_no:10}  seats:{seat_no:11}  seats:{seat_no:12}  seats:{seat_no:13}  seats:{seat_no:14}  seats:{seat_no:15}  seats:{seat_no:16}  seats:{seat_no:17}  seats:{seat_no:18}  seats:{seat_no:19}  seats:{seat_no:20}  seats:{seat_no:21}  seats:{seat_no:22}  seats:{seat_no:23}  seats:{seat_no:24}  seats:{seat_no:25}  seats:{seat_no:26}  seats:{seat_no:27}  seats:{seat_no:28}  seats:{seat_no:29}  seats:{seat_no:30}  seats:{seat_no:31}  seats:{seat_no:32}  seats:{seat_no:33}  seats:{seat_no:34}  seats:{seat_no:35}  seats:{seat_no:36}  seats:{seat_no:37}  seats:{seat_no:38}  seats:{seat_no:39}  seats:{seat_no:40}  seats:{seat_no:41}  seats:{seat_no:42}  seats:{seat_no:43}  seats:{seat_no:44}  seats:{seat_no:45}  seats:{seat_no:46}  seats:{seat_no:47}  seats:{seat_no:48}  seats:{seat_no:49}  seats:{seat_no:50}}


 User successfully removed from train


 User allocation modified successfully


 Recevied allocation data for section A:  train_number:1  section_details:{seats:{seat_no:1}  seats:{seat_no:2}  seats:{seat_no:3}  seats:{seat_no:4}  seats:{seat_no:5}  seats:{seat_no:6}  seats:{seat_no:7}  seats:{seat_no:8}  seats:{seat_no:9}  seats:{seat_no:10}  seats:{seat_no:11}  seats:{seat_no:12}  seats:{seat_no:13}  seats:{seat_no:14}  seats:{seat_no:15}  seats:{seat_no:16}  seats:{seat_no:17}  seats:{seat_no:18}  seats:{seat_no:19}  seats:{seat_no:20}  seats:{seat_no:21}  seats:{seat_no:22}  seats:{seat_no:23}  seats:{seat_no:24}  seats:{seat_no:25}  seats:{seat_no:26}  seats:{seat_no:27}  seats:{seat_no:28}  seats:{seat_no:29}  seats:{seat_no:30}  seats:{seat_no:31}  seats:{seat_no:32}  seats:{seat_no:33}  seats:{seat_no:34}  seats:{seat_no:35}  seats:{seat_no:36}  seats:{seat_no:37}  seats:{seat_no:38}  seats:{seat_no:39}  seats:{seat_no:40}  seats:{seat_no:41}  seats:{seat_no:42}  seats:{seat_no:43}  seats:{seat_no:44}  seats:{seat_no:45}  seats:{seat_no:46}  seats:{seat_no:47}  seats:{seat_no:48}  seats:{seat_no:49}  seats:{seat_no:50}}


 Recevied allocation data for section B:  train_number:1  section_details:{name:B  seats:{seat_no:1  is_allocated:true  allocated_to:{first_name:"Manoj"  last_name:"Kumar"  email:"manunero1730@gmail.com"}}  seats:{seat_no:2}  seats:{seat_no:3}  seats:{seat_no:4}  seats:{seat_no:5}  seats:{seat_no:6}  seats:{seat_no:7}  seats:{seat_no:8}  seats:{seat_no:9}  seats:{seat_no:10  is_allocated:true  allocated_to:{first_name:"Manoj"  last_name:"Kumar"  email:"manunero1730@gmail.com"}}  seats:{seat_no:11}  seats:{seat_no:12}  seats:{seat_no:13}  seats:{seat_no:14}  seats:{seat_no:15}  seats:{seat_no:16}  seats:{seat_no:17}  seats:{seat_no:18}  seats:{seat_no:19}  seats:{seat_no:20}  seats:{seat_no:21}  seats:{seat_no:22}  seats:{seat_no:23}  seats:{seat_no:24}  seats:{seat_no:25}  seats:{seat_no:26}  seats:{seat_no:27}  seats:{seat_no:28}  seats:{seat_no:29}  seats:{seat_no:30}  seats:{seat_no:31}  seats:{seat_no:32}  seats:{seat_no:33}  seats:{seat_no:34}  seats:{seat_no:35}  seats:{seat_no:36}  seats:{seat_no:37}  seats:{seat_no:38}  seats:{seat_no:39}  seats:{seat_no:40}  seats:{seat_no:41}  seats:{seat_no:42}  seats:{seat_no:43}  seats:{seat_no:44}  seats:{seat_no:45}  seats:{seat_no:46}  seats:{seat_no:47}  seats:{seat_no:48}  seats:{seat_no:49}  seats:{seat_no:50}}
ticketmaster % 

```




