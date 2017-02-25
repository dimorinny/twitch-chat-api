package twitchchat

const (
	host = "irc.chat.twitch.tv"
)

type (
	Configuration struct {
		Host     string
		Nickname string
		Oauth    string
		Channel  string
	}
)

func NewConfiguration(nickname, oauth, channel string) (configuration *Configuration) {
	configuration = &Configuration{
		Host:     host,
		Nickname: nickname,
		Oauth:    oauth,
		Channel:  channel,
	}
	return
}
