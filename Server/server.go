package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	gRPC "github.com/seve0039/Distributed-Auction-System.git/proto"
	"google.golang.org/grpc"
)

type Server struct {
	gRPC.UnimplementedAuctionServiceServer
	name          string
	port          string
	mapOfBidders  map[int64]string
	auctionIsOpen bool
}

var serverName = flag.String("name", "default", "Server's name")
var port = flag.String("port", "5400", "Server port")
var currentHighestBid int64

func main() {
	flag.Parse()
	launchServer()
}

func launchServer() {
	list, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", *port))
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", *port, err)
	}

	auctionServer := grpc.NewServer()
	server := &Server{
		name:          *serverName,
		port:          *port,
		mapOfBidders:  make(map[int64]string),
		auctionIsOpen: true,
		//participants: make(map[string]gRPC.AuctionService_BroadcastServer),
	}
	gRPC.RegisterAuctionServiceServer(auctionServer, server)
	//gRPC.RegisterAuctionServiceServer(auctionServer, server)
	log.Printf("NEW SESSION: Server %s: Listening at %v\n", *serverName, list.Addr())

	if err := auctionServer.Serve(list); err != nil {
		log.Fatalf("failed to serve %v", err)
	}
}

func (s *Server) Bid(context context.Context, bidAmount *gRPC.BidAmount) (*gRPC.Ack, error) {
	if s.auctionIsOpen == true {
		higher := isHigherThanCurrentBid(bidAmount.Amount)
		if higher {
			s.mapOfBidders[bidAmount.Amount] = bidAmount.Name
			return &gRPC.Ack{Acknowledgement: "Success"}, nil
		} else {
			return &gRPC.Ack{Acknowledgement: "Fail"}, nil
		}
	}
	return &gRPC.Ack{Acknowledgement: "Auction is not yet open!"}, nil
}

func isHigherThanCurrentBid(bidAmount int64) (isHigher bool) {
	if currentHighestBid < bidAmount {
		currentHighestBid = bidAmount
		return true
	} else {
		return false
	}

}

/*func (s *Server) Result(context context.Context, empty google.protobuf.Empty) (*gRPC.HighestBid){
	return &gRPC.HighestBid{HighestBid: currentHighestBid, Name: s.mapOfBidders[currentHighestBid]}
}*/
