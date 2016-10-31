package store

type Store interface {
	GetUidsByChatId(chatId string) []string
	GetPushAddrByUid(uid string) string
	JoinChat(version string, chatId string, uid string) error
}
