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

## Usage

Ferstly, you need to create chat object like this:

```
twitch := twitchchat.NewChat(config)
```

For receiving chat status (like connected, disconnected, new message) you can use 2 ways:

### Using callbacks api

```
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
```

### Using channels api

```
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
```

For more complicated usage example see [sample code](https://github.com/dimorinny/twitch-chat-api/blob/master/sample/main.go).
