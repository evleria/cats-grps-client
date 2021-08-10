package main

import (
	"context"
	"github.com/evleria/cats-app/protocol/pb"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial(":6000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal("cannot connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewCatsServiceClient(conn)

	printAllCats(c)
}

func printAllCats(client pb.CatsServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	stream, err := client.GetAllCats(ctx, &empty.Empty{})
	if err != nil {
		log.Fatalf("cannot get all cats: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Fatalf("cannot get message: %v", err)
		}

		log.Printf("got cat: %s", resp.Cat.Name)
	}
	log.Println("fetched all cats successfully")
}
