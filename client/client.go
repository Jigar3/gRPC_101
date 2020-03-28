package main

import (
	"context"
	"fmt"
	"log"

	githubp "github.com/jigar3/grpc/githubpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Gihhub client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Can't connect to server: %v", err)
	}
	defer cc.Close()

	c := githubp.NewGithubServiceClient(cc)
	getF(c)
}

func getF(c githubp.GithubServiceClient) {
	fmt.Println("Starting getF Function")

	req := &githubp.FollowerRequest{
		GithubUsername: "jigar3",
	}

	res, err := c.GetFollowers(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)
	}

	fmt.Println(len(res.FollowerList))
	log.Printf("The Result is: %v", res.FollowerList)
}
