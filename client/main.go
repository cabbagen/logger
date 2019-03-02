package main

import (
	"context"
	"log"
	"time"
	
	"google.golang.org/grpc"
	pb "logger/proto"
)

const address = "localhost:50051"

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	
	defer conn.Close()
	
	c := pb.NewLoggerServiceClient(conn)
	
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	
	defer cancel()
	
	appid := "1"
	message := "hello world"
	category := "console"
	
	loggerRequest := pb.LoggerRequest {
		Appid: &appid,
		Message: &message,
		Category: &category,
	}
	
	r, err := c.AllLog(ctx, &loggerRequest)
	
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	
	log.Printf("Greeting %s", *(r.Message))
}