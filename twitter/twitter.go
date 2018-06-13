package twitter

import (
	"fmt"
	"log"

	"github.com/ChimeraCoder/anaconda"
	"github.com/joho/godotenv"
)

type TwitterToken struct {
	AccessToken       string
	AccessTokenSecret string
}

func SetupTwitter() {
	var myEnv map[string]string
	myEnv, err := godotenv.Read()
	if err != nil {
		log.Fatal(err)
	}
	anaconda.SetConsumerKey(myEnv["TWITTER_CONSUMER_KEY"])
	anaconda.SetConsumerSecret(myEnv["TWITTER_CONSUMER_SECRET"])
}

func Tweet(token TwitterToken, text string) {
	api := anaconda.NewTwitterApi(token.AccessToken, token.AccessTokenSecret)
	tweet, err := api.PostTweet(text, nil)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(tweet.Text)
}

func DirectMessageWithId(token TwitterToken, text string, toUser int64) {
	api := anaconda.NewTwitterApi(token.AccessToken, token.AccessTokenSecret)
	dm, err := api.PostDMToUserId(text, toUser)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(dm.Text)
}

func DirectMessageWithName(token TwitterToken, text string, toUser string) {
	api := anaconda.NewTwitterApi(token.AccessToken, token.AccessTokenSecret)
	dm, err := api.PostDMToScreenName(text, toUser)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(dm.Text)
}
