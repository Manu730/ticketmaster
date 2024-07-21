package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"

	pb "ticketmaster/protogen/golang"

	"google.golang.org/grpc"
)

// Generate sample train data
func generateTrainData() []*pb.Train {
	var trains []*pb.Train
	tripDetail := &pb.Trip{Source: "London", Destination: "France"}
	var sections []*pb.Section
	for i := 1; i <= 2; i++ {
		var seats []*pb.Seat
		for j := 1; j <= 50; j++ {
			seat := &pb.Seat{SeatNo: int64(j), IsAllocated: false}
			seats = append(seats, seat)
		}
		section := &pb.Section{Name: pb.SectionOrder(i - 1), Seats: seats}
		sections = append(sections, section)
	}

	train := &pb.Train{Number: 1, TripDetail: tripDetail, Sections: sections}
	trains = append(trains, train)
	return trains
}

// Server object
type ticketMaster struct {
	pb.UnimplementedTicketMasterServer
	trains     []*pb.Train
	lock       sync.Mutex
	bookingMap map[string]*pb.Receipt
}

// Generate a new server object
func NewTicketMaster() *ticketMaster {
	tMaster := &ticketMaster{}
	tMaster.trains = generateTrainData()
	tMaster.bookingMap = make(map[string]*pb.Receipt)
	return tMaster
}

// API to book ticket for a given user
func (tm *ticketMaster) BookTicket(ctx context.Context, in *pb.Receipt) (*pb.BookTicketOutput, error) {
	fmt.Println("\n\n Received request for ticekt booking for user: ", in.User.Email)
	if in.FromStation == "" {
		return nil, errors.New("source station cannot be empty")
	} else if in.ToStation == "" {
		return nil, errors.New("destination station cannot be empty")
	} else if in.User == nil {
		return nil, errors.New("user cannot be empty")
	} else if in.Cost == nil {
		return nil, errors.New("ticket cost cannot be empty")
	}
	tm.lock.Lock()
	defer tm.lock.Unlock()
	isSeatAllocated := false
	passenger := &pb.User{FirstName: in.User.FirstName, LastName: in.User.LastName, Email: in.User.Email}
	var currTrip *pb.Trip
	for _, train := range tm.trains {
		if train.TripDetail.Source == in.FromStation && train.TripDetail.Destination == in.ToStation {
			for _, sec := range train.Sections {
				for _, se := range sec.Seats {
					if !se.IsAllocated {
						se.AllocatedTo = passenger
						se.IsAllocated = true
						isSeatAllocated = true
						currTrip = train.TripDetail
						break
					}
				}
			}
		}
	}
	if !isSeatAllocated {
		return nil, errors.New("Train is full")
	}

	ticketprice := &pb.Price{Price: in.Cost.Price, Currency: in.Cost.Currency}
	receipt := &pb.Receipt{TrainNumber: 1, FromStation: in.FromStation, ToStation: in.ToStation, User: passenger, Cost: ticketprice}
	tm.bookingMap[in.User.Email] = receipt
	bookingoutput := &pb.BookTicketOutput{TrainNumber: 1, TripDetail: currTrip, Traveller: in.User}
	return bookingoutput, nil
}

// API to get receipt for a given user
func (tm *ticketMaster) GetReceipt(ctx context.Context, in *pb.UserTrainInput) (*pb.Receipt, error) {
	fmt.Println("\n\n Received request for receipt for user: ", in.User.Email)
	tm.lock.Lock()
	defer tm.lock.Unlock()
	if _, ok := tm.bookingMap[in.User.Email]; !ok {
		return nil, errors.New("No booking for the user")
	}
	return tm.bookingMap[in.User.Email], nil
}

// API to show allocations for a given section of the train
func (tm *ticketMaster) ShowAllocations(ctx context.Context, in *pb.ShowAllocationInput) (*pb.ShowAllocationOutput, error) {
	fmt.Println("\n\n Received request to show ticket allocations")
	tm.lock.Lock()
	defer tm.lock.Unlock()
	for _, train := range tm.trains {
		if train.Number == in.TrainNumber {
			for _, sec := range train.Sections {
				if sec.Name == in.Order {
					return &pb.ShowAllocationOutput{TrainNumber: in.TrainNumber, SectionDetails: sec}, nil
				}
			}
		}
	}
	return nil, errors.New("section not found")
}

// API to remove given user from train
func (tm *ticketMaster) RemoveUser(ctx context.Context, in *pb.UserTrainInput) (*pb.Empty, error) {
	fmt.Println("\n\n Received request to remove user: ", in.User.Email)
	tm.lock.Lock()
	defer tm.lock.Unlock()
	for _, train := range tm.trains {
		if train.Number == in.TrainNumber {
			for _, sec := range train.Sections {
				for _, se := range sec.Seats {
					if se.IsAllocated && se.AllocatedTo.Email == in.User.Email {
						se.AllocatedTo = nil
						se.IsAllocated = false
						break
					}
				}
			}
		}
	}
	delete(tm.bookingMap, in.User.Email)
	return &pb.Empty{}, nil
}

// API to modify user allocation
func (tm *ticketMaster) ModifyUserAllocation(ctx context.Context, in *pb.UserAllocModifyInput) (*pb.Empty, error) {
	fmt.Println("\n\n Received request to modify allocation for user: ", in.UserTrainDetails.User.Email)
	tm.lock.Lock()
	defer tm.lock.Unlock()
	passenger := &pb.User{FirstName: in.UserTrainDetails.User.FirstName, LastName: in.UserTrainDetails.User.LastName, Email: in.UserTrainDetails.User.Email}
	for _, train := range tm.trains {
		if train.Number == in.UserTrainDetails.TrainNumber {
			for _, sec := range train.Sections {
				if sec.Name == in.NewSection {
					for _, se := range sec.Seats {
						if se.SeatNo == in.NewSeatno && !se.IsAllocated {
							se.AllocatedTo = passenger
							se.IsAllocated = true
						}
					}
				} else {
					for _, se := range sec.Seats {
						if se.IsAllocated && se.AllocatedTo.Email == in.UserTrainDetails.User.Email {
							se.IsAllocated = false
							se.AllocatedTo = nil
						}
					}
				}
			}
		}
	}
	return &pb.Empty{}, nil
}

func main() {
	server := grpc.NewServer()
	ticketMas := NewTicketMaster()
	pb.RegisterTicketMasterServer(server, ticketMas)
	tl, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(fmt.Println("Error starting tcp listener on port 8080", err))
	} else {
		fmt.Println("listening on port 8080")
	}
	server.Serve(tl)
}
