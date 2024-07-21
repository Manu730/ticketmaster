package main

import (
	"context"
	"errors"
	"log"
	"net"
	"testing"

	pb "ticketmaster/protogen/golang"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func server(ctx context.Context) (pb.TicketMasterClient, func()) {
	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)

	baseServer := grpc.NewServer()
	pb.RegisterTicketMasterServer(baseServer, NewTicketMaster())
	go func() {
		if err := baseServer.Serve(lis); err != nil {
			log.Printf("error serving server: %v", err)
		}
	}()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error connecting to server: %v", err)
	}

	closer := func() {
		err := lis.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		baseServer.Stop()
	}

	client := pb.NewTicketMasterClient(conn)

	return client, closer
}

func TestTicketMaster_BookTicket(t *testing.T) {
	ctx := context.Background()
	client, closer := server(ctx)
	defer closer()
	type expectation struct {
		out *pb.BookTicketOutput
		err error
	}
	tests := map[string]struct {
		in       *pb.Receipt
		expected expectation
	}{
		"success": {
			in:       &pb.Receipt{TrainNumber: 1, FromStation: "London", ToStation: "France", User: &pb.User{FirstName: "Manoj", LastName: "Kumar", Email: "manunero1730@gmail.com"}, Cost: &pb.Price{Price: 20, Currency: "Dollar"}},
			expected: expectation{out: &pb.BookTicketOutput{TrainNumber: 1, TripDetail: &pb.Trip{Source: "London", Destination: "France"}, Traveller: &pb.User{FirstName: "Manoj", LastName: "Kumar", Email: "manunero1730@gmail.com"}}, err: nil},
		},

		"source_station_empty": {
			in:       &pb.Receipt{TrainNumber: 1, FromStation: "", ToStation: "France", User: &pb.User{FirstName: "Manoj", LastName: "Kumar", Email: "manunero1730@gmail.com"}, Cost: &pb.Price{Price: 20, Currency: "Dollar"}},
			expected: expectation{out: nil, err: errors.New("rpc error: code = Unknown desc = source station cannot be empty")},
		},

		"destination_station_empty": {
			in:       &pb.Receipt{TrainNumber: 1, FromStation: "London", ToStation: "", User: &pb.User{FirstName: "Manoj", LastName: "Kumar", Email: "manunero1730@gmail.com"}, Cost: &pb.Price{Price: 20, Currency: "Dollar"}},
			expected: expectation{out: nil, err: errors.New("rpc error: code = Unknown desc = destination station cannot be empty")},
		},

		"user_empty": {
			in:       &pb.Receipt{TrainNumber: 1, FromStation: "London", ToStation: "France", User: nil, Cost: &pb.Price{Price: 20, Currency: "Dollar"}},
			expected: expectation{out: nil, err: errors.New("rpc error: code = Unknown desc = user cannot be empty")},
		},
		"cost_empty": {
			in:       &pb.Receipt{TrainNumber: 1, FromStation: "London", ToStation: "France", User: &pb.User{FirstName: "Manoj", LastName: "Kumar", Email: "manunero1730@gmail.com"}, Cost: nil},
			expected: expectation{out: nil, err: errors.New("rpc error: code = Unknown desc = ticket cost cannot be empty")},
		},
	}
	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			out, err := client.BookTicket(ctx, tt.in)
			if err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Err -> \nWant: %q\nGot: %q\n", tt.expected.err.Error(), err.Error())
				}
			} else {
				if tt.expected.out.Traveller.Email != out.Traveller.Email {
					t.Errorf("Err -> \nWant: %q\nGot: %q\n", tt.expected.out.Traveller.Email, out.Traveller.Email)
				}

			}
		})
	}
}

func TestTicketMaster_GetReceipt(t *testing.T) {
	ctx := context.Background()
	client, closer := server(ctx)
	test_traveller := &pb.User{FirstName: "Manoj", LastName: "Kumar", Email: "manunero1730@gmail.com"}
	test_receipt := &pb.Receipt{TrainNumber: 1, FromStation: "London", ToStation: "France", User: test_traveller, Cost: &pb.Price{Price: 20, Currency: "Dollar"}}
	defer closer()
	_, err := client.BookTicket(ctx, test_receipt)
	if err != nil {
		t.Errorf("Error booking ticket: %q", err.Error())
	}
	type expectation struct {
		out *pb.Receipt
		err error
	}
	tests := map[string]struct {
		in       *pb.UserTrainInput
		expected expectation
	}{
		"success": {
			in:       &pb.UserTrainInput{TrainNumber: 1, User: test_traveller},
			expected: expectation{out: test_receipt, err: nil},
		},
	}
	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			o, _ := client.GetReceipt(ctx, tt.in)
			if tt.expected.out.User.Email != o.User.Email {
				t.Errorf("Err -> \nWant: %q\nGot: %q\n", tt.expected.out.User.Email, o.User.Email)
			}

		})
	}

}

func TestTicketMaster_ShowAllocations(t *testing.T) {
	ctx := context.Background()
	client, closer := server(ctx)
	defer closer()
	var seats []*pb.Seat
	for j := 1; j <= 50; j++ {
		seat := &pb.Seat{SeatNo: int64(j), IsAllocated: false}
		seats = append(seats, seat)
	}
	section := &pb.Section{Name: pb.SectionOrder(0), Seats: seats}

	type expectation struct {
		out *pb.ShowAllocationOutput
		err error
	}
	tests := map[string]struct {
		in       *pb.ShowAllocationInput
		expected expectation
	}{
		"success": {
			in:       &pb.ShowAllocationInput{TrainNumber: 1, Order: 0},
			expected: expectation{out: &pb.ShowAllocationOutput{TrainNumber: 1, SectionDetails: section}, err: nil},
		},
	}
	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			s, _ := client.ShowAllocations(ctx, tt.in)
			if len(tt.expected.out.SectionDetails.Seats) != len(s.SectionDetails.Seats) {
				t.Errorf("Err -> \nWant: %q\nGot: %q\n", len(tt.expected.out.SectionDetails.Seats), len(s.SectionDetails.Seats))
			}
		})
	}
}

func TestTicketMaster_RemoveUser(t *testing.T) {

}

func TestTicketMaster_ModifyUserAllocation(t *testing.T) {

}
