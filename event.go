package twitchchat

type (
	Connected    func()
	Disconnected func()
	NewMessage   func(message string)
)
