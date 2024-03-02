package transport

const (
	TgBotDirectChatID int64 = 6406834985
)

// Commands:
const (
	EditChannelsTopics = "edit_channel_topics"
	ListChannelsTopics = "list_channels_topics"
	MyChannels         = "my_channels"

	Moderate        = "moderate"
	ModerateDetails = "moderate_details"
	PostNow         = "post_now"
	PostLater       = "post_later"
	RejectAd        = "reject_ad"
)

type BotState int

const (
	StateStart BotState = iota
	StateEditTopics
)
