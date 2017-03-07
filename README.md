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
```

### Using channels api

```
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
```

For more complicated usage example see [sample code](https://github.com/dimorinny/twitch-chat-api/blob/master/sample/main.go).
