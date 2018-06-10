package twitter

import (
	"fmt"
	"log"

	"github.com/ChimeraCoder/anaconda"
	"github.com/joho/godotenv"
)

type TwitterToken struct {
	accessToken       string
	accessTokenSecret string
}

func SetupTwitter() {
	var myEnv map[string]string
	myEnv, err := godotenv.Read("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	anaconda.SetConsumerKey(myEnv["TWITTER_CONSUMER_KEY"])
	anaconda.SetConsumerSecret(myEnv["TWITTER_CONSUMER_SECRET"])
}

func Tweet(token TwitterToken, text string) {
	api := anaconda.NewTwitterApi(token.accessToken, token.accessTokenSecret)
	tweet, err := api.PostTweet(text, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tweet.Text)
}

func DirectMessageWithId(token TwitterToken, text string, toUser int64) {
	api := anaconda.NewTwitterApi(token.accessToken, token.accessTokenSecret)
	dm, err := api.PostDMToUserId(text, toUser)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dm.Text)
}

func DirectMessageWithName(token TwitterToken, text string, toUser string) {
	api := anaconda.NewTwitterApi(token.accessToken, token.accessTokenSecret)
	dm, err := api.PostDMToScreenName(text, toUser)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dm.Text)
}
