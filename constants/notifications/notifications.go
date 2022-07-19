package notifications

const (
	USAGE_NAME               = "usage: /name [name]"
	USAGE_CREATE             = "usage: /create [channelName]"
	USAGE_JOIN               = "usage: /join [channelName]"
	CLIENT_CHANGED_NAME      = "client changed their name"
	CLIENT_JOIN_CHANNEL      = "client joined"
	CLIENT_CREATED_CHANNEL   = "client created channel"
	CLIENT_CONNECTION_CLOSED = "client closed the connection"
	CLIENT_LIST_CHANNELS     = "client listed channels"
	INVALID_REQUEST          = "client request not supported"
)
