package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dghubble/go-twitter/twitter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type config struct {
	Key    string `json:"key"`
	Secret string `json:"secret"`
	User   struct {
		ID   int64  `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"user"`
	Limit int64 `json:"limit"`
}

type truncUser struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Public        bool   `json:"public"`
	LastTweetID   int64  `json:"last_tweet_id"`
	LastTweetText string `json:"last_tweet_text"`
}

func main() {
	b, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	var conf config
	err = json.Unmarshal(b, &conf)
	if err != nil {
		log.Fatal(err)
	}

	config := &clientcredentials.Config{
		ClientID:     conf.Key,
		ClientSecret: conf.Secret,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}

	cli := config.Client(oauth2.NoContext)
	twt := twitter.NewClient(cli)

	fid, res, err := twt.Friends.IDs(&twitter.FriendIDParams{
		UserID:     conf.User.ID,
		ScreenName: conf.User.Name,
	})
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		log.Fatal(res.Status)
	}

	users, res, err := twt.Users.Lookup(&twitter.UserLookupParams{
		UserID: fid.IDs[:conf.Limit],
	})
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		log.Fatal(res.Status)
	}

	var tu []truncUser
	for _, user := range users {
		if user.Status == nil {
			tu = append(tu, truncUser{
				ID:     user.ID,
				Name:   user.ScreenName,
				Public: !user.Protected,
			})
		} else {
			tu = append(tu, truncUser{
				ID:            user.ID,
				Name:          user.ScreenName,
				Public:        !user.Protected,
				LastTweetID:   user.Status.ID,
				LastTweetText: user.Status.Text,
			})
		}
	}

	b, err = json.MarshalIndent(tu, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("verified.json", b, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
