package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	gRPC "github.com/seve0039/Distributed-Auction-System.git/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var user = generateRandomString(5)
var clientsName = flag.String("name", user, "Sender's name")
var serverPort = flag.String("server", "5400", "Tcp server")

var server gRPC.AuctionServiceClient
var ServerConn *grpc.ClientConn

func main() {
	flag.Parse()
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

	server = gRPC.NewAuctionServiceClient(conn)
	ServerConn = conn
	fmt.Println("Connected to server")
}

func sendBid(amount int64) { //Make a bid
	ack, _ := server.Bid(context.Background(), &gRPC.BidAmount{
		Id: 1, Amount: amount, Name: *clientsName,
	})
	fmt.Println(ack.Acknowledgement)

}

func handleCommand() { //Handle commands from user input via the terminal
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("--- Please make your bid ---")
	fmt.Println("Write 'help' for options")
	for {
		fmt.Print(" ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		input = strings.TrimSpace(input)
		if input == "status" {
			fmt.Println("The highest bid for now is : [COMMING SOON]") //TODO: Get results
		} else if input == "help" {
			fmt.Println("-- To see the currnet highest bid write 'status' and press enter")
			fmt.Println("-- To place a bid write the amount you want to bid as a <number> and press enter")
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
			return
		}

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

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
