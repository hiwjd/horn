package store

type Store interface {
	GetUidsByChatId(chatId string) []string
	GetPushAddrByUid(uid string) string
	JoinChat(mid string, chatId string, uid string, role string) error
	CreateChat(chatId string, cid string, creator string, staffId string) error
	GetStaffsByCompany(cid string) []string
	GetChatsByUid(uid string) []string
	LeaveChat(mid string, chatId string, uid string) error
}
