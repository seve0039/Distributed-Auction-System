package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	gRPC "github.com/seve0039/Distributed-Auction-System.git/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

var user = ""
var clientsName = flag.String("name", user, "Sender's name")
var serverPort = flag.String("server", "5400", "Tcp server")

var server gRPC.AuctionServiceClient
var ServerConn *grpc.ClientConn

func main() {
	flag.Parse()
	createClientName()
	connectToServer()
	sendStreamConnection()
	handleCommand()

	for {

	}

}

func connectToServer() {
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial(fmt.Sprintf(":%s", *serverPort), opts...)
	if err != nil {
		log.Fatalf("Fail to Dial : %v", err)
		fmt.Println("Failed to connect to server")
	}
	fmt.Println("Connected to server")
	server = gRPC.NewAuctionServiceClient(conn)
	ServerConn = conn

}

func sendBid(amount int64) { //Make a bid
	ack, _ := server.Bid(context.Background(), &gRPC.BidAmount{
		Id: 1, Amount: amount, Name: *clientsName,
	})
	fmt.Println(ack.Acknowledgement)
}

func getResult() { //Get the result of the auction
	result, _ := server.Result(context.Background(), &emptypb.Empty{})
	fmt.Println("The current highest bid is:", result.HighestBid, "kr. by:", result.HighestBidderName)
}

func handleCommand() { //Handle commands from user input via the terminal
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Write 'help' for options")
	fmt.Println("--- Please make your bid ---")
	for {
		fmt.Print(" ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		input = strings.TrimSpace(input)
		if input == "help" {

			fmt.Println("--- You have the following options ---")
			fmt.Println("- 'status'	to see the currnet highest bid")
			fmt.Println("- <amount>	to place a bid with an amount of your choice")

		} else if input == "status" {
			getResult()

		} else {

			var bid int64
			fmt.Sscan(input, &bid)
			sendBid(bid)

		}

	}
}

func listenForResult(stream gRPC.AuctionService_BroadcastToAllClient) { //Listen for results from the server
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Println("Failed to receive broadcast: ", err)

			time.Sleep(40 * time.Second)
			fmt.Println("Attempting reconection to server")
			stream, err := server.BroadcastToAll(context.Background())
			if err != nil {
				log.Fatalf("Error while creating stream: %v", err)
			}
			stream.Send(&gRPC.StreamConnection{StreamName: *clientsName})
		}
		fmt.Println("The auction has ended")

		fmt.Println(msg.StreamName)
	}
}

func sendStreamConnection() {
	stream, err := server.BroadcastToAll(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
	}
	stream.Send(&gRPC.StreamConnection{StreamName: *clientsName})
	go listenForResult(stream)
}

func createClientName() {
	if *clientsName == "" {
		fmt.Print("Enter your name: ")
		reader := bufio.NewReader(os.Stdin)
		name, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Errer reading name: %v", err)
		}
		*clientsName = strings.TrimSpace(name)
	}

}
