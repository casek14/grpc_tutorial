package main

import (
	pb "first_grpc/customer"
	"context"
	"strings"
	"net"
	"log"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

//server is used to implement customer.CustomerServer
type server struct {
	savedCustomers []*pb.CustomerRequest

}

//CreateCustomer creates a new customer
func (s *server) CreateCustomer(ctx context.Context, in *pb.CustomerRequest)(*pb.CustomerResponse, error){
	s.savedCustomers = append(s.savedCustomers, in)
	return &pb.CustomerResponse{Id: in.Id, Sucess: true}, nil
}


func (s *server) GetCustomers(filter *pb.CustomerFilter, stream pb.Customer_GetCustomersServer) error{
	for _, customer := range s.savedCustomers{
		if filter.Keyword != ""{
			if !strings.Contains(customer.Name, filter.Keyword){
				continue
			}
		}
		err := stream.Send(customer)
		if err != nil {
			return err
		}


		}
		return nil
	}


func main() {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v ",err)
	}

	//create new gRPC server
	s := grpc.NewServer()
	pb.RegisterCustomerServer(s, &server{})
	s.Serve(lis)
	}
