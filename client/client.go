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
	if err != nil {
		fmt.Println("\n\n Error connecting to server: ", err)
		panic(err)
	}
	client := pb.NewTicketMasterClient(conn)
	ctx := context.Background()
	traveller := &pb.User{FirstName: "Manoj", LastName: "Kumar", Email: "manunero1730@gmail.com"}
	payment := &pb.Price{Price: 20, Currency: "Dollar"}

	//Book ticket for a given user
	bookTicketInput := &pb.Receipt{TrainNumber: 1, FromStation: "London", ToStation: "France", User: traveller, Cost: payment}
	bookingresp, er := client.BookTicket(context.Background(), bookTicketInput)
	if er != nil {
		fmt.Println("\n\n Error Booking ticket: ", er)
		panic(er)

	}
	fmt.Println("\n\n Recevied booking response: ", bookingresp)

	//Get receipt for the booked user
	receipt, e := client.GetReceipt(ctx, &pb.UserTrainInput{TrainNumber: 1, User: traveller})
	if e != nil {
		fmt.Println("\n\n Error fetching data: ", e)
	}
	fmt.Println("\n\n Recevied receipt: ", receipt)

	//Show train allocation post ticket booking for section A
	sectionAdata, err1 := client.ShowAllocations(ctx, &pb.ShowAllocationInput{TrainNumber: 1, Order: 0})
	if err1 != nil {
		fmt.Println("\n\n Error fetching allocation data: ", err1)
	} else {
		fmt.Println("\n\n Recevied allocation data for section A: ", sectionAdata)
	}

	//Show train allocation post ticket booking for section B
	sectionBdata, err2 := client.ShowAllocations(ctx, &pb.ShowAllocationInput{TrainNumber: 1, Order: 1})
	if err2 != nil {
		fmt.Println("\n\n Error fetching allocation data: ", err2)
	} else {
		fmt.Println("\n\n Recevied allocation data for section B: ", sectionBdata)
	}

	//Remove User from train
	_, err3 := client.RemoveUser(ctx, &pb.UserTrainInput{TrainNumber: 1, User: &pb.User{FirstName: "Manoj", LastName: "Kumar", Email: "manunero1730@gmail.com"}})
	if err3 != nil {
		fmt.Println("\n\n Error removing user from train: ", err3)
	} else {
		fmt.Println("\n\n User successfully removed from train")
	}

	//Book ticket again for user
	client.BookTicket(context.Background(), bookTicketInput)
	_, err4 := client.ModifyUserAllocation(ctx, &pb.UserAllocModifyInput{UserTrainDetails: &pb.UserTrainInput{TrainNumber: 1, User: &pb.User{FirstName: "Manoj", LastName: "Kumar", Email: "manunero1730@gmail.com"}}, NewSection: 1, NewSeatno: 10})
	if err4 != nil {
		fmt.Println("\n\n Error updating user allocation: ", err4)
	} else {
		fmt.Println("\n\n User allocation modified successfully")
	}

	//Show train allocation post user modification for section A
	sectionAdata, err1 = client.ShowAllocations(ctx, &pb.ShowAllocationInput{TrainNumber: 1, Order: 0})
	if err1 != nil {
		fmt.Println("\n\n Error fetching allocation data: ", err1)
	} else {
		fmt.Println("\n\n Recevied allocation data for section A: ", sectionAdata)
	}

	//Show train allocation post user modification for section B
	sectionBdata, err2 = client.ShowAllocations(ctx, &pb.ShowAllocationInput{TrainNumber: 1, Order: 1})
	if err2 != nil {
		fmt.Println("\n\n Error fetching allocation data: ", err2)
	} else {
		fmt.Println("\n\n Recevied allocation data for section B: ", sectionBdata)
	}

}
