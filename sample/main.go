package main

import (
	"fmt"
	"github.com/dimorinny/twitch-chat-api"
)

const (
	nickname = "<your-nickname>"
	channel  = "<streamer's-nickname>"

	// Your chat scoped Twitch auth token.
	// For example: oauth:v5lw2a2mnv18q3a0iey40fwewdparm
	oauth = "oauth:<your-token>"
)

var (
	config *twitchchat.Configuration
)

func initConfiguration() {
	config = twitchchat.NewConfiguration(nickname, oauth, channel)
}

func init() {
	initConfiguration()
}

func main() {
	twitch := twitchchat.NewChat(config)

	runWithCallbacks(twitch)
}

func runWithChannels(twitch *twitchchat.Chat) {
	stop := make(chan struct{})
	defer close(stop)

	disconnected := make(chan struct{})
	connected := make(chan struct{})
	message := make(chan string)

	go func() {
		for {
			select {
			case <-disconnected:
				fmt.Println("Disconnected")
				stop <- struct{}{}
			case <-connected:
				fmt.Println("Connected")
			case newMessage := <-message:
				fmt.Println(newMessage)
			}
		}
	}()

	if err := twitch.ConnectWithChannels(connected, disconnected, message); err != nil {
		return
	}

	<-stop
}

func runWithCallbacks(twitch *twitchchat.Chat) {
	stop := make(chan struct{})
	defer close(stop)

	err := twitch.ConnectWithCallbacks(
		func() {
			fmt.Println("Connected")
		},
		func() {
			fmt.Println("Disconnected")
			stop <- struct{}{}
		},
		func(message string) {
			fmt.Println(message)
		},
	)

	if err != nil {
		return
	}

	<-stop
}
