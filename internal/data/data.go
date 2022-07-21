package data

type ChannelData struct {
	Realm      string
	User       string
	Allocation string
}

type TurnData struct {
	SentP int64
	RecvP int64
	SentB int64
	RecvB int64
}
