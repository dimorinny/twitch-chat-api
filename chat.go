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
	message chan<- string,
) error {
	closeChannels := func() {
		close(connected)
		close(disconnected)
		close(message)
	}

	connectedCallback := func() {
		connected <- struct{}{}
	}

	disconnectedCallback := func() {
		disconnected <- struct{}{}
		closeChannels()
	}

	newMessageCallback := func(newMessage string) {
		message <- newMessage
	}

	err := c.ConnectWithCallbacks(
		connectedCallback,
		disconnectedCallback,
		newMessageCallback,
	)

	if err != nil {
		closeChannels()
		return err
	}

	return nil
}

func (c *Chat) ConnectWithCallbacks(
	connected Connected,
	disconnected Disconnected,
	message NewMessage,
) error {
	c.connection.HandleFunc(connectedEvent, func(conn *irc.Conn, line *irc.Line) {
		connected()
		c.connection.Join("#" + c.config.Channel)
	})

	c.connection.HandleFunc(disconnectedEvent, func(conn *irc.Conn, line *irc.Line) {
		disconnected()
	})

	c.connection.HandleFunc(newMessageEvent, func(conn *irc.Conn, line *irc.Line) {
		if len(line.Args) > 1 {
			message(line.Args[1])
		}
	})

	if err := c.connection.Connect(); err != nil {
		return err
	}

	return nil
}
