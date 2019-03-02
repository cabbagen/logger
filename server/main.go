package main

import (
	"log"
	"net"
	"google.golang.org/grpc"
	"logger/server/server"
	pb "logger/proto"
)

const port = ":50051"

func main() {
	listen, error := net.Listen("tcp", port)

	if error != nil {
		log.Fatalf("failed to listen: %v\n", error)
	}

	s := grpc.NewServer()

	pb.RegisterLoggerServiceServer(s, server.NewRPCServer())

	log.Printf("start server on %s\n", port)

	if error := s.Serve(listen); error != nil {
		log.Fatalf("failed to server: %v\n", error)
	}
}










