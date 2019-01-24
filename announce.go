package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	yaml "gopkg.in/yaml.v2"

	"github.com/nlopes/slack"
)

func main() {
	fileLocation := os.Args[1]

	announcements := loadAllAnnouncements(fileLocation)

	for _, announcement := range announcements {
		err := sendToTwitter(announcement)
		if err != nil {
			fmt.Printf("Cound not send to twitter")
		}
		err = postAndPinToSlack(announcement)
		if err != nil {
			fmt.Printf("Cound not send to twitter")
		}
	}
}

// AnnouncementData contains all the information needed to send a message to
// a specific endpoint like twitter, slack, etc
type AnnouncementData struct {
	Title   string
	Type    string
	Message string
	// Slack fields
	ChannelID string `yaml:"channel_id"`
	Token     string
	DoPin     bool `yaml:"do_pin"`

	// Twitter fields
	ConsumerKey    string `yaml:"consumer_key"`
	ConsumerSecret string `yaml:"consumer_secret"`
	AccessKey      string `yaml:"access_key"`
	AccessSecret   string `yaml:"access_secret"`
}

// loadAllAnnouncements will look at a folder and attempt to load all the yaml
// files containing annoucnement data
func loadAllAnnouncements(location string) []AnnouncementData {
	announcements := []AnnouncementData{}
	err := filepath.Walk(location, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ".yaml" {
			return nil
		}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Printf("Could not load file from %v \n", path)
			return nil
		}
		announcement := AnnouncementData{}
		err = yaml.Unmarshal([]byte(data), &announcement)
		if err != nil {
			log.Printf("Could not unmarshall file at %v: %v", path, err)
			return nil
		}
		fmt.Printf("--- t:\n%v\n\n", announcement)
		announcements = append(announcements, announcement)
		return nil
	})
	if err != nil {
		fmt.Printf("Could not walk file tree at %v \n: %v", location, err)
		return nil
	}
	fmt.Printf("--- t:\n%v\n\n", announcements)
	return announcements
}

func sendToTwitter(announcement AnnouncementData) error {
	config := oauth1.NewConfig(announcement.ConsumerKey, announcement.ConsumerSecret)
	token := oauth1.NewToken(announcement.AccessKey, announcement.AccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	// Send message to twitter
	client := twitter.NewClient(httpClient)
	tweet, resp, err := client.Statuses.Update(announcement.Message, nil)
	if err != nil {
		fmt.Printf("Error sending announcement text to twitter: %v", err)
		return err
	}

	fmt.Printf("tweet: %v", tweet)
	fmt.Printf("resp: %v", resp)
	fmt.Printf("err: %v", err)
	return nil
}

func postAndPinToSlack(announcement AnnouncementData) error {
	api := slack.New(announcement.Token)

	channelID, timestamp, err := api.PostMessage(announcement.ChannelID, announcement.Message, slack.NewPostMessageParameters())
	if err != nil {
		fmt.Printf("Error sending message to slack: %s\n", err)
		return err
	}

	msgRef := slack.NewRefToMessage(channelID, timestamp)

	// Add message pin to channel
	if announcement.DoPin == true {
		if err = api.AddPin(channelID, msgRef); err != nil {
			fmt.Printf("Error adding pin to slack: %s\n", err)
			return err
		}
	}

	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	return nil
}
