## High level Twitch chat abstraction

This library provide high level wrapper above Twitch (IRC) API, that allow you analyze stream comments and building Twitch bots.

## Installation

```
go get github.com/dimorinny/twitch-chat-api
```

## Configuration

Create configuration object:

```
config = twitchchat.NewConfiguration(
	nickname,
	oauth,
	channel,
)
```

You can quickly get a oauth token for your account with this [helpful page](http://twitchapps.com/tmi/).
For more information you should read official Twitch IRC [documentation](https://github.com/justintv/Twitch-API/blob/master/IRC.md).
