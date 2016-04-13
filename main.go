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
	"github.com/rk/go-cron"
)

type returnObject struct {
	Path   string `json:"path"`
	Secret string `json:"secret"`
}

func (r *returnObject) JSON() []byte {
	v, _ := json.Marshal(r)
	return v
}

func init() {
	cron.NewDailyJob(11, 0, 0, func(time.Time) {
		if !state {
			PostTweet()
		}
		state = false
	})
}

func main() {
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

	Tweet, err = api.PostTweet("Testing...", v)
	if err != nil {
		log.Println("We have an error now ", err)
	}
}

var (
	adjs        = []string{"autumn", "hidden", "bitter", "misty", "silent", "empty", "dry", "dark", "summer", "icy", "delicate", "quiet", "white", "cool", "spring", "winter", "patient", "twilight", "dawn", "crimson", "wispy", "weathered", "blue", "billowing", "broken", "cold", "damp", "falling", "frosty", "green", "long", "late", "lingering", "bold", "little", "morning", "muddy", "old", "red", "rough", "still", "small", "sparkling", "throbbing", "shy", "wandering", "withered", "wild", "black", "young", "holy", "solitary", "fragrant", "aged", "snowy", "proud", "floral", "restless", "divine", "polished", "ancient", "purple", "lively", "nameless"}
	nouns       = []string{"waterfall", "river", "breeze", "moon", "rain", "wind", "sea", "morning", "snow", "lake", "sunset", "pine", "shadow", "leaf", "dawn", "glitter", "forest", "hill", "cloud", "meadow", "sun", "glade", "bird", "brook", "butterfly", "bush", "dew", "dust", "field", "fire", "flower", "firefly", "feather", "grass", "haze", "mountain", "night", "pond", "darkness", "snowflake", "silence", "sound", "sky", "shape", "surf", "thunder", "violet", "water", "wildflower", "wave", "water", "resonance", "sun", "wood", "dream", "cherry", "tree", "fog", "frost", "voice", "paper", "frog", "smoke", "star"}
	Tweet       anaconda.Tweet
	secret      = "secret"
	state       = false
	stateObject = returnObject{Path: "/start", Secret: "startSecret"}
	errorObject = returnObject{Path: "error", Secret: "Secret not matched"}
)

const (
	consumerKey    = "<from twitter>"
	consumerSecret = "<from twitter>"
	accessToken    = "<from twitter>"
	accessSecret   = "<from twitter>"
	macAddress     = "<from twitter>"
)
