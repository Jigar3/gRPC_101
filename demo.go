package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

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

func getFollower(username string) (err error) {
	fmt.Printf("The Username for which list of followers was requested is %s \n", username)

	client := &http.Client{}
	page := 1
	url := fmt.Sprintf("https://api.github.com/users/%s/followers?per_page=100&page=%d", username, page)

	first_req, _ := http.NewRequest("GET", url, nil)

	first_resp, err := client.Do(first_req)

	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}

	defer first_resp.Body.Close()
	first_resp_body, _ := ioutil.ReadAll(first_resp.Body)

	var apiResponse APIResponse

	json.Unmarshal(first_resp_body, &apiResponse)

	for {
		page = page + 1
		if len(apiResponse)%100 == 0 {
			url := fmt.Sprintf("https://api.github.com/users/%s/followers?per_page=100&page=%d", username, page)
			r, _ := http.NewRequest("GET", url, nil)

			resp, _ := client.Do(r)

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

	for _, n := range apiResponse {
		fmt.Println(n.Login)
	}

	return nil
}

func main() {
	getFollower("mubaris")
}
