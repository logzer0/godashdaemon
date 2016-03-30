package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

type returnObject struct {
	Path   string `json:"path"`
	Secret string `json:"secret"`
}

func (r *returnObject) JSON() []byte {
	v, _ := json.Marshal(r)
	return v
}

var (
	Tweet              anaconda.Tweet
	haveWeTweetedToday = true
	secret             = "secret"
	state              = true
	stateObject        = returnObject{Path: "/start", Secret: "startSecret"}
	errorObject        = returnObject{Path: "error", Secret: "Secret not matched"}
)

func main() {
	loc, _ := time.LoadLocation("America/Chicago")
	ticker := time.NewTicker(time.Hour * 24)
	go func() {
		for _ = range ticker.C {
			currentTime := time.Now().In(loc)
			cutOfftime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 10, 0, 0, 0, loc)
			if time.Now().After(cutOfftime) {
				if state {
					//All is good. Do nothing
					if !haveWeTweetedToday {
						haveWeTweetedToday = true
						//fmt.Println("I'm a good boy")
					} else {
						DeleteTweet()
						//fmt.Println("This came after the cutoff time, so we are deleting the earlier tweet")
					}
				} else {
					//This means we failed. So, the failure case goes here
					if !haveWeTweetedToday {
						tweetTime = time.Now()
						PostTweet()
						haveWeTweetedToday = true
						//fmt.Println("Bitch we are tweeting now")
					} else {
						//fmt.Println("It's not set yet, but since we tweeted already we won't be posting the tweet")
					}
				}
			}
			fmt.Println("EOD. Resetting stuff")
			state = false
			haveWeTweetedToday = false
			Tweet = nil
		}
	}()

	http.HandleFunc(stateObject.Path, handler)
	http.ListenAndServe(":9580", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.FormValue(secret) == stateObject.Secret && r.URL.Path == stateObject.Path {
		fmt.Println("State is now set")
		state = true
		stateObject.Path, stateObject.Secret = "/"+generateName(), generateName()
		w.Write(stateObject.JSON())
		http.HandleFunc(stateObject.Path, handler)
	} else {
		w.Write(errorObject.JSON())
	}
}

func generateName() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%s-%s-%d", adjs[r.Intn(len(adjs))], nouns[r.Intn(len(nouns))], r.Intn(9999-1000)+1000)
}

func PostTweet() {
	fmt.Println("Now posting tweet")
	var err error
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessSecret)

	v := url.Values{}

	rand.Seed(time.Now().UTC().UnixNano())
	//Tweet, err = api.PostTweet(statements[rand.Intn(len(statements))], v)
	Tweet, err = api.PostTweet("Testing...", v)
	if err != nil {
		log.Println("We have an error now ", err)
	}
}

func DeleteTweet() {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessSecret)

	_, err := api.DeleteTweet(Tweet.Id, true)
	if err != nil {
		log.Println("We have an error now ", err)
	}
}

var (
	adjs       = []string{"autumn", "hidden", "bitter", "misty", "silent", "empty", "dry", "dark", "summer", "icy", "delicate", "quiet", "white", "cool", "spring", "winter", "patient", "twilight", "dawn", "crimson", "wispy", "weathered", "blue", "billowing", "broken", "cold", "damp", "falling", "frosty", "green", "long", "late", "lingering", "bold", "little", "morning", "muddy", "old", "red", "rough", "still", "small", "sparkling", "throbbing", "shy", "wandering", "withered", "wild", "black", "young", "holy", "solitary", "fragrant", "aged", "snowy", "proud", "floral", "restless", "divine", "polished", "ancient", "purple", "lively", "nameless"}
	nouns      = []string{"waterfall", "river", "breeze", "moon", "rain", "wind", "sea", "morning", "snow", "lake", "sunset", "pine", "shadow", "leaf", "dawn", "glitter", "forest", "hill", "cloud", "meadow", "sun", "glade", "bird", "brook", "butterfly", "bush", "dew", "dust", "field", "fire", "flower", "firefly", "feather", "grass", "haze", "mountain", "night", "pond", "darkness", "snowflake", "silence", "sound", "sky", "shape", "surf", "thunder", "violet", "water", "wildflower", "wave", "water", "resonance", "sun", "wood", "dream", "cherry", "tree", "fog", "frost", "voice", "paper", "frog", "smoke", "star"}
	statements = []string{
		"PHP is the coolest language",
		"I love Javascript",
		"I love working with Internet Explorer",
		"Functional programming is overrated",
		"LISP is for losers",
		"Windows phone is the next big thing",
	}
)

const (
	consumerKey    = "KZpZPlN9tAMJPQuEwkWW9eJJW"
	consumerSecret = "xpPuFe2a8NA8uXASIrFCiJvpxw4EYImq18y8lDH3MpNjfelEQz"
	accessToken    = "41533711-lsMZ9lzAy336cQ7Mf8qlk0YXxTBXmFLGSky5VyVW7"
	accessSecret   = "yNywgviBcfVeNCFbKB6j6kXslutwHfc1FDn8Du8pqG4vf"
	macAddress     = " 74:c2:46:84:f5:ee"
)
