package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"bufio"

	gRPC "github.com/seve0039/Distributed-Auction-System.git/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// var clientsName = flag.String("name", "Placeholder", "Sender's name")
var serverPort = flag.String("server", "5400", "Tcp server")

var server gRPC.AuctionServiceClient
var ServerConn *grpc.ClientConn

func main() {
	flag.Parse()
	connectToServer()

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

	server = gRPC.NewAuctionServiceClient(conn)
	ServerConn = conn
	fmt.Println("Connected to server")
}

func sendMessage(amount int64) {
	_, _ = server.Bid(context.Background(), &gRPC.BidAmount{
		Id: 1, Amount: amount, Name: "John Doe",
	})
}

func handleCommand() { //Handle commands from user input via the terminal
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("--- Please make your bid ---")
		for {
			input, err := reader.ReadString('\n')
			if err != nil {
				log.Fatal(err)
		}
		if input == "status" {
			fmt.Println("The highest bid for now is :") //TODO: Get results
		} else {
			
		}

	}
}
