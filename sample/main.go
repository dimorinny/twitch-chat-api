package main

import (
	"fmt"
	"github.com/dimorinny/twitch-chat-api"
	irc "github.com/fluffle/goirc/client"
)

const (
	nickname = "<your-nickname>"
	channel  = "<streamer's-nickname>"

	// Your chat scoped Twitch auth token.
	// For example: oauth:v5lw2a2mnv18q3a0iey40fwewdparm
	oauth = "oauth:<your-token>"
)

var (
	config     *twitchchat.Configuration
	connection *irc.Conn
)

func initConfiguration() {
	config = twitchchat.NewConfiguration(nickname, oauth, channel)
}

func initTwitchConnection() {
	ircConfig := irc.NewConfig(config.Nickname)
	ircConfig.Server = config.Host
	ircConfig.Pass = config.Oauth

	connection = irc.Client(ircConfig)
}

func init() {
	initConfiguration()
	initTwitchConnection()
}

func main() {
	twitch := twitchchat.NewChat(
		config,
		connection,
	)

	runWithCallbacks(twitch)
}

func runWithChannels(twitch *twitchchat.Chat) {
	disconnected := make(chan struct{})
	connected := make(chan struct{})
	errStream := make(chan error)
	message := make(chan string)

	go func() {
		for {
			select {
			case <-disconnected:
				fmt.Println("Disconnected")
			case <-connected:
				fmt.Println("Connected")
			case err := <-errStream:
				fmt.Println(err)
			case newMessage := <-message:
				fmt.Println(newMessage)
			}
		}
	}()

	twitch.ConnectWithChannels(connected, disconnected, errStream, message)
}

func runWithCallbacks(twitch *twitchchat.Chat) {
	twitch.ConnectWithCallbacks(
		func() {
			fmt.Println("Connected")
		},
		func() {
			fmt.Println("Disconnected")
		},
		func(err error) {
			fmt.Println(err)
		},
		func(message string) {
			fmt.Println(message)
		},
	)
}
