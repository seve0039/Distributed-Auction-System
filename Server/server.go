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

func main() {
	flag.Parse()
	createLogFile()
	//go endAuction()
	go launchServer()

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

func launchServer() {
	list, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", *port))
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", *port, err)
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

func (s *Server) Bid(context context.Context, bidAmount *gRPC.BidAmount) (*gRPC.Ack, error) {
	if s.auctionIsOpen {

		higher := isHigherThanCurrentBid(bidAmount.Amount)
		if higher {
			s.mapOfBidders[bidAmount.Amount] = bidAmount.Name
			return &gRPC.Ack{Acknowledgement: "Success"}, nil
		} else {
			return &gRPC.Ack{Acknowledgement: "Fail"}, nil
		}

	}
	return &gRPC.Ack{Acknowledgement: "Auction is not open yet!"}, nil
}

func (s *Server) Result(context context.Context, empty *emptypb.Empty) (*gRPC.HighestBid, error) {
	return &gRPC.HighestBid{HighestBid: currentHighestBid, HighestBidderName: s.mapOfBidders[currentHighestBid]}, nil
}

func (s *Server) BroadcastToAll(stream gRPC.AuctionService_BroadcastToAllServer) error {
	for {

		in, err := stream.Recv()
		if err != nil {
			return err
		}
		s.participants[in.StreamName] = stream
		fmt.Println("New participant: ", in.StreamName)
	}
}

// This function sends the result to all participants
func sendResult(message string) {
	//This is the line to run When auction ends
	//<sendResult(fmt.Sprintf("Highest bid is %d by %s", bidAmount.Amount, bidAmount.Name))>

	for _, participant := range server.participants {
		participant.Send(&gRPC.StreamConnection{StreamName: message})
	}

}

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

func createLogFile() {
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)
}

/*func (s *Server) Result(context context.Context, empty google.protobuf.Empty) (*gRPC.HighestBid){
	return &gRPC.HighestBid{HighestBid: currentHighestBid, Name: s.mapOfBidders[currentHighestBid]}
}*/
