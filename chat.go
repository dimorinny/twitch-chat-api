package twitchchat

import (
	irc "github.com/fluffle/goirc/client"
	"github.com/fluffle/goirc/state"
)

const (
	connectedEvent    = "connected"
	disconnectedEvent = "disconnected"
	newMessageEvent   = "privmsg"
)

type (
	Chat struct {
		config     *Configuration
		connection *irc.Conn
	}
)

func NewChat(config *Configuration) *Chat {
	ircConfig := irc.NewConfig(config.Nickname)
	ircConfig.Server = config.Host
	ircConfig.Pass = config.Oauth

	return &Chat{
		config:     config,
		connection: irc.Client(ircConfig),
	}
}

func NewChatWithIrc(config *Configuration, ircConfig *irc.Config) *Chat {
	ircConfig.Me = &state.Nick{Nick: config.Nickname}
	ircConfig.Server = config.Host
	ircConfig.Pass = config.Oauth

	return &Chat{
		config:     config,
		connection: irc.Client(ircConfig),
	}
}

func (c *Chat) ConnectWithChannels(
	connected, disconnected chan<- struct{},
	errorStream chan<- error,
	message chan<- string,
) {
	connectedCallback := func() {
		connected <- struct{}{}
	}

	disconnectedCallback := func() {
		disconnected <- struct{}{}
	}

	errorCallback := func(err error) {
		errorStream <- err
	}

	newMessageCallback := func(newMessage string) {
		message <- newMessage
	}

	c.ConnectWithCallbacks(
		connectedCallback,
		disconnectedCallback,
		errorCallback,
		newMessageCallback,
	)
}

func (c *Chat) ConnectWithCallbacks(
	connected Connected,
	disconnected Disconnected,
	error Error,
	message NewMessage,
) {
	quit := make(chan struct{})

	c.connection.HandleFunc(connectedEvent, func(conn *irc.Conn, line *irc.Line) {
		connected()
		c.connection.Join("#" + c.config.Channel)
	})

	c.connection.HandleFunc(disconnectedEvent, func(conn *irc.Conn, line *irc.Line) {
		disconnected()
		quit <- struct{}{}
	})

	c.connection.HandleFunc(newMessageEvent, func(conn *irc.Conn, line *irc.Line) {
		if len(line.Args) > 1 {
			message(line.Args[1])
		}
	})

	if err := c.connection.Connect(); err != nil {
		error(err)
		return
	}

	<-quit
}
