package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	gRPC "github.com/seve0039/Distributed-Auction-System.git/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	gRPC.UnimplementedAuctionServiceServer
	name          string
	port          string
	participants  map[string]gRPC.AuctionService_BroadcastToAllServer
	mapOfBidders  map[int64]string
	auctionIsOpen bool
}

var auctionIsOpen bool = true
var server *Server
var serverName = flag.String("name", "default", "Server's name")
var port = flag.String("port", "5400", "Server port")
var currentHighestBid int64
var ports [3]string = [3]string{*port, "5401", "5402"} //Ports to connect to in case of server crash
var portCounter = 1                                    // Counts the number of used ports

func main() {
	flag.Parse()
	createLogFile()
	go launchServer(*port)

	for {
		if !auctionIsOpen {
			auctionIsOpen = true
			fmt.Println("I get here")
			sendResult(fmt.Sprintf("Highest bid was %d by %s", currentHighestBid, server.mapOfBidders[currentHighestBid]))
			time.Sleep(1 * time.Second)
			os.Exit(1)

		}
	}
}

// Launching the server and starts the auction
func launchServer(_ string) {
	list, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", *port))
	if err != nil {
		if portCounter != len(ports) {
			//log.Fatalf("Failed to listen on port %s: %v", *port, err)
			log.Printf("Server %s: Trying to find another port", *serverName)
			port := ports[portCounter]
			portCounter++      // Counts the number of used ports
			launchServer(port) // Relaunch server with new port
		} else {
			log.Fatalf("Server %s: Failed to find available port", *serverName)
		}
	}

	auctionServer := grpc.NewServer()
	server = &Server{
		name:          *serverName,
		port:          *port,
		participants:  make(map[string]gRPC.AuctionService_BroadcastToAllServer),
		mapOfBidders:  make(map[int64]string),
		auctionIsOpen: true,
	}
	gRPC.RegisterAuctionServiceServer(auctionServer, server)
	//gRPC.RegisterAuctionServiceServer(auctionServer, server)
	log.Printf("NEW SESSION: Server %s: Listening at %v\n", *serverName, list.Addr())
	if err := auctionServer.Serve(list); err != nil {
		log.Fatalf("failed to serve %v", err)
	}

}

// From the .proto file. Handles a bid from a client during the auction and returns the acknowledgement (Success or fail)
func (s *Server) Bid(context context.Context, bidAmount *gRPC.BidAmount) (*gRPC.Ack, error) {
	if s.auctionIsOpen {
		higher := isHigherThanCurrentBid(bidAmount.Amount)
		if higher {
			s.mapOfBidders[bidAmount.Amount] = bidAmount.Name
			log.Println("Participant", bidAmount.Name, "is now the highest bidder with:", bidAmount.Amount)
			fmt.Println("Participant", bidAmount.Name, "is now the highest bidder with:", bidAmount.Amount)
			return &gRPC.Ack{Acknowledgement: "Success: You are now the highest bidder"}, nil
		} else {
			log.Println("Participant", bidAmount.Name, "got rejected with the bid:", bidAmount.Amount)
			fmt.Println("Participant", bidAmount.Name, "got rejected with the bid:", bidAmount.Amount)
			return &gRPC.Ack{Acknowledgement: "Fail: Your bid was too low"}, nil
		}

	}
	return &gRPC.Ack{Acknowledgement: "Auction is not open yet!"}, nil
}

func (s *Server) Result(context context.Context, empty *emptypb.Empty) (*gRPC.HighestBid, error) {
	return &gRPC.HighestBid{HighestBid: currentHighestBid, HighestBidderName: s.mapOfBidders[currentHighestBid]}, nil
}

// From the .proto file. Broadcasting to all clients
func (s *Server) BroadcastToAll(stream gRPC.AuctionService_BroadcastToAllServer) error {
	for {
		in, err := stream.Recv()
		if err != nil {
			return err
		}
		s.participants[in.StreamName] = stream
		log.Println("New participant: ", in.StreamName)
		fmt.Println("New participant: ", in.StreamName)
	}
}

// This function sends the result to all participants
func sendResult(message string) {
	for _, participant := range server.participants {
		participant.Send(&gRPC.StreamConnection{StreamName: message})
	}

}

// Returns true if the bid is higher than the highest bid
func isHigherThanCurrentBid(bidAmount int64) (isHigher bool) {
	if currentHighestBid < bidAmount {
		currentHighestBid = bidAmount
		return true
	} else {
		return false
	}

}
func endAuction() {
	time.Sleep(15 * time.Second)
	auctionIsOpen = false
	fmt.Println("Auction is now closed")
}

// Creates and connects to the log.txt file
func createLogFile() {
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)
}


