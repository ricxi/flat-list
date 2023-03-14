package main

import (
	"context"
	"fmt"
	"log"

	"githug.com/ricxi/flat-list/token/activation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cc, err := grpc.Dial(":5003", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("client could not connect: %s\n", err)
	}
	defer cc.Close()

	c := activation.NewTokenServiceClient(cc)

	req := activation.Request{UserId: "fakeittillyoumakeit"}
	res, err := c.CreateTokenForUser(context.Background(), &req)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Response: ", res)
}
