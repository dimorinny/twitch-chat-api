package twitchchat

type (
	Connected    func()
	Disconnected func()
	Error        func(err error)
	NewMessage   func(message string)
)
