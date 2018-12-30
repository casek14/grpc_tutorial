package main

import (
	pb "first_grpc/customer"
	"context"
	"log"
	"io"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func CreateCustomer(client pb.CustomerClient, customer *pb.CustomerRequest){
	resp, err := client.CreateCustomer(context.Background(),customer)
	if err != nil {
		log.Fatalf("Could not create Customer: %v \n", err)
	}

	if resp.Sucess {
		log.Printf("A new customer has been added with id: %d \n", resp.Id)
	}
}

func GetCustomers(client pb.CustomerClient, filter  *pb.CustomerFilter){

	//calling streaming API
	stream, err := client.GetCustomers(context.Background(),filter)
	if err != nil {
		log.Fatalf("Error on get customers: %v\n", err)
	}

	for {
		customer, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("%v.GetCustomers(_) = _,%v\n", client, err)
		}

		log.Printf("Customer: %v\n", customer)
	}
}

func main() {

	conn, err := grpc.Dial(address,grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}

	defer conn.Close()
	//Create new customer client
	client := pb.NewCustomerClient(conn)

	customer1 := &pb.CustomerRequest{
		Id: 101,
		Name: "Karel Novak",
		Email: "novak12@seznam.cz",
		Phone: "111222333",
		Addresses: []*pb.CustomerRequest_Address{
			&pb.CustomerRequest_Address{
				Street: "Sladkovskeho",
				City: "Pardubice",
				State: "Czech Republic",
				Zip: "53374",
				IsShippingAddress: true,
			},
			&pb.CustomerRequest_Address{
				Street: "Pod vilami",
				City: "Praha",
				State: "Czech Republic",
				Zip: "50002",
				IsShippingAddress: false,
			},
		},
	}

	// create new customer
	CreateCustomer(client, customer1)

	customer2 := &pb.CustomerRequest{
		Id: 102,
		Name: "Iren Rose",
		Email: "i2rose@gmail.com",
		Phone: "725-432-565",
		Addresses: []*pb.CustomerRequest_Address{
			&pb.CustomerRequest_Address{
				Street: "1 st mission street",
				City: "San Francisco",
				State: "CA",
				Zip: "94105",
				IsShippingAddress: true,
			},
		},
	}

	CreateCustomer(client,customer2)

	//Filter with an empty keyword
	filter := &pb.CustomerFilter{Keyword:""}
	GetCustomers(client,filter)

	}
