package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
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
	secret      = "secret"
	state       bool
	stateObject = returnObject{Path: "/start", Secret: "startSecret"}
	errorObject = returnObject{Path: "error", Secret: "Secret not matched"}
)

func main() {
	loc, _ := time.LoadLocation("America/Chicago")
	ticker := time.NewTicker(time.Millisecond * 5000)
	go func() {
		for _ = range ticker.C {
			currentTime := time.Now().In(loc)
			cutOfftime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 10, 0, 0, 0, loc)
			if time.Now().After(cutOfftime) {
				if state {
					//All is good. Do nothing
				} else {
					//This means we failed. So, the failure case goes here
				}
				state = false
			}
		}
	}()

	http.HandleFunc(stateObject.Path, handler)
	http.ListenAndServe(":9580", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.FormValue(secret) == stateObject.Secret && r.URL.Path == stateObject.Path {
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

var (
	adjs  = []string{"autumn", "hidden", "bitter", "misty", "silent", "empty", "dry", "dark", "summer", "icy", "delicate", "quiet", "white", "cool", "spring", "winter", "patient", "twilight", "dawn", "crimson", "wispy", "weathered", "blue", "billowing", "broken", "cold", "damp", "falling", "frosty", "green", "long", "late", "lingering", "bold", "little", "morning", "muddy", "old", "red", "rough", "still", "small", "sparkling", "throbbing", "shy", "wandering", "withered", "wild", "black", "young", "holy", "solitary", "fragrant", "aged", "snowy", "proud", "floral", "restless", "divine", "polished", "ancient", "purple", "lively", "nameless"}
	nouns = []string{"waterfall", "river", "breeze", "moon", "rain", "wind", "sea", "morning", "snow", "lake", "sunset", "pine", "shadow", "leaf", "dawn", "glitter", "forest", "hill", "cloud", "meadow", "sun", "glade", "bird", "brook", "butterfly", "bush", "dew", "dust", "field", "fire", "flower", "firefly", "feather", "grass", "haze", "mountain", "night", "pond", "darkness", "snowflake", "silence", "sound", "sky", "shape", "surf", "thunder", "violet", "water", "wildflower", "wave", "water", "resonance", "sun", "wood", "dream", "cherry", "tree", "fog", "frost", "voice", "paper", "frog", "smoke", "star"}
)
