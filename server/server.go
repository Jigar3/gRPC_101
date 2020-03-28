package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	githubp "github.com/jigar3/grpc/githubpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

type APIResponse []struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

func (*server) GetFollowers(ctx context.Context, req *githubp.FollowerRequest) (res *githubp.FollowerResponse, err error) {
	fmt.Printf("The username for which list of followers was requested is %s \n", req.GithubUsername)

	client := &http.Client{}
	page := 1
	url := fmt.Sprintf("https://api.github.com/users/%s/followers?per_page=100&page=%d", req.GithubUsername, page)

	first_req, _ := http.NewRequest("GET", url, nil)

	first_resp, err := client.Do(first_req)

	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return &githubp.FollowerResponse{}, err
	}

	defer first_resp.Body.Close()
	first_resp_body, _ := ioutil.ReadAll(first_resp.Body)

	var apiResponse APIResponse

	json.Unmarshal(first_resp_body, &apiResponse)

	for {
		page = page + 1
		if len(apiResponse)%100 == 0 {
			url := fmt.Sprintf("https://api.github.com/users/%s/followers?per_page=100&page=%d", req.GithubUsername, page)
			r, _ := http.NewRequest("GET", url, nil)

			resp, err := client.Do(r)
			if err != nil {
				fmt.Println("Errored when sending request to the server")
				return &githubp.FollowerResponse{}, err
			}

			defer resp.Body.Close()
			r_body, _ := ioutil.ReadAll(resp.Body)

			var temp APIResponse

			json.Unmarshal(r_body, &temp)
			apiResponse = append(apiResponse, temp...)

			if len(temp) == 0 {
				break
			}
		} else {
			break
		}
	}

	var followerList []string
	for _, n := range apiResponse {
		followerList = append(followerList, n.Login)
	}

	return &githubp.FollowerResponse{FollowerList: followerList}, nil
}

func main() {
	fmt.Println("Github Service")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	githubp.RegisterGithubServiceServer(s, &server{})

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
