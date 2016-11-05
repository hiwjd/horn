package store

type Store interface {
	GetUidsByChatId(chatId string) []string
	GetPushAddrByUid(uid string) string
	JoinChat(mid string, chatId string, uid string, role string) error
	CreateChat(chatId string, gid string, creator string, kfid int) error
}
