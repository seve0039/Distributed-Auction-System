package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	gRPC "github.com/seve0039/Distributed-Auction-System.git/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	gRPC.UnimplementedAuctionServiceServer
	name string
	port string
	//participants map[string]gRPC.AuctionService_BroadcastServer
}

var serverName = flag.String("name", "default", "Server's name")
var port = flag.String("port", "5400", "Server port")

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
		name: *serverName,
		port: *port,
		//participants: make(map[string]gRPC.AuctionService_BroadcastServer),
	}

	gRPC.RegisterAuctionServiceServer(auctionServer, server)
	log.Printf("NEW SESSION: Server %s: Listening at %v\n", *serverName, list.Addr())

	if err := auctionServer.Serve(list); err != nil {
		log.Fatalf("failed to serve %v", err)
	}
}

func (s *Server) Bid(context context.Context, criticalSectionRequest *gRPC.BidAmount) (*emptypb.Empty, error) {

}
