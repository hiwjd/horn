package store

type Store interface {
	GetUidsByChatId(chatId string) []string
	GetPushAddrByUid(uid string) string
}
