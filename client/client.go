package main

import (
	"context"
	"fmt"
	pb "ticketmaster/protogen/golang"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := pb.NewTicketMasterClient(conn)

	traveller := &pb.User{FirstName: "Manoj", LastName: "Kumar", Email: "manunero1730@gmail.com"}
	payment := &pb.Price{Price: 20, Currency: "Dollar"}
	bookTicketInput := &pb.Receipt{TrainNumber: 1, FromStation: "London", ToStation: "France", User: traveller, Cost: payment}

	bookingresp, er := client.BookTicket(context.Background(), bookTicketInput)
	if er != nil {
		fmt.Println("Error fetching data: ", er)
	}
	fmt.Println("Recevied booking resp: ", bookingresp)

	allocationdata, err := client.ShowAllocations(context.Background(), &pb.ShowAllocationInput{TrainNumber: 1, Order: 0})
	if err != nil {
		fmt.Println("Error fetching data: ", err)
	}
	fmt.Println("Recevied data: ", allocationdata)

	receipt, e := client.GetReceipt(context.Background(), &pb.UserTrainInput{TrainNumber: 1, User: traveller})
	if e != nil {
		fmt.Println("Error fetching data: ", e)
	}
	fmt.Println("Recevied receipt: ", receipt)
}
