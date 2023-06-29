package wsManager

var Manager = NewClientManager()
var GroupChannelsManager = NewGroupChannelsManager()

func DoInit() {
	ToClientChan = make(chan *clientInfo, 2048)
	PingTimer()
	SubRedis()
}
