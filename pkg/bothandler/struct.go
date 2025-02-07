package bothandler

// We could use discordbot's handler system, but we have own wrappers here to make
// multiplatform bots work.

type MessageHandler func() string
type CatchallHandler func(string) string

type MessagePlatform interface {
	Send(string)
	ProcessMessages()
	Close()
	ChannelMessageSend(channel string, message string) error
}

type AddMessagePlatform func(MessagePlatform)

var Handlers = map[string]MessageHandler{}
var CatchallHandlers = []CatchallHandler{}
var AddMessagePlatforms = []AddMessagePlatform{}
var ActiveMessagePlatforms = []MessagePlatform{}

func RegisterMessageHandler(m string, h MessageHandler) {
	Handlers[m] = h
}

func RegisterCatchallHandler(h CatchallHandler) {
	CatchallHandlers = append(CatchallHandlers, h)
}

// Yes, weird. All the modules register themselves,
// Then all the platforms (discord, telegram, slack) calls these
// to tell the modules these platforms are available for them
// to call asynchronously.
// We need these to get the moduls to run themselves asychronously
// apart from the main messaging platform event loop
func RegisterPlatformRegisteration(h AddMessagePlatform) {
	AddMessagePlatforms = append(AddMessagePlatforms, h)
}

func RegisterMessagePlatform(m MessagePlatform) {
	for _, v := range AddMessagePlatforms {
		v(m)
	}
	ActiveMessagePlatforms = append(ActiveMessagePlatforms, m)
}

func RegisterPassiveMessagePlatform(m MessagePlatform) {
	ActiveMessagePlatforms = append(ActiveMessagePlatforms, m)
}

func Shutdown() {
	for _, v := range ActiveMessagePlatforms {
		v.Close()
	}
}

func ChannelMessageSend(channelId string, message string) error {
	for _, v := range ActiveMessagePlatforms {
		err := v.ChannelMessageSend(channelId, message)
		if err != nil {
			return err
		}
	}
	return nil
}
