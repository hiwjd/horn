package consumer

import (
	"encoding/json"
)

// 消息来源方信息
type From struct {
	Id   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

// 对话信息
type Chat struct {
	Id string `db:"id" json:"id"`
}

// 消息基本信息 会匿名组合到具体的消息里
type Message struct {
	Type string         `db:"type" json:"type"` // 消息的类型 text, file, image, event
	T    map[string]int `db:"t" json:"t"`       // 0:客户端发出时间戳 1:入队列时间戳 2:分发时间戳
	Mid  string         `db:"mid" json:"mid"`   // 消息ID
	From From           `db:"from" json:"from"` // 消息发送方信息
}

// 普通消息
type MessageText struct {
	Message
	Chat Chat   `db:"chat" json:"chat"` // 对话信息
	Text string `db:"text" json:"text"`
}

// 图片数据
type Image struct {
	Src    string `db:"src" json:"src"`
	Width  int    `db:"width" json:"width"`
	Height int    `db:"height" json:"height"`
	Size   int    `db:"size" json:"size"`
}

// 图片消息
type MessageImage struct {
	Message
	Chat  Chat  `db:"chat" json:"chat"` // 对话信息
	Image Image `db:"image" json:"image"`
}

// 文件数据
type File struct {
	Src  string `db:"src" json:"src"`
	Name string `db:"name" json:"name"`
	Size int    `db:"size" json:"size"`
}

// 文件消息
type MessageFile struct {
	Message
	Chat Chat `db:"chat" json:"chat"` // 对话信息
	File File `db:"file" json:"file"`
}

// 请求对话数据
type EventRequestChat struct {
	Chat Chat     `db:"chat" json:"chat"` // 对话信息
	Uids []string `db:"uids" json:"uids"` // 邀请加入对话的uid数组
}

// 请求对话消息
type MessageEventRequestChat struct {
	Message
	Event EventRequestChat `db:"event" json:"event"`
}

// 加入对话数据
type EventJoinChat struct {
	Chat Chat `db:"chat" json:"chat"` // 对话信息
}

// 加入对话消息
type MessageEventJoinChat struct {
	Message
	Event EventJoinChat `db:"event" json:"event"`
}

type Message2Pusher struct {
	Type string      `db:"type" json:"type"`
	To   []string    `db:"to" json:"to"`
	Data interface{} `db:"data" json:"data"`
}

// 访问数据
type MessageViewPage struct {
	TrackId string `db:"track_id" json:"track_id"`
	Uid     string `db:"uid" json:"uid"`
	Fp      string `db:"fb" json:"fp"`
	Gid     string `db:"gid" json:"gid"`
	Url     string `db:"url" json:"url"`
	Title   string `db:"title" json:"title"`
	Referer string `db:"referer" json:"referer"`
	Os      string `db:"os" json:"os"`
	Browser string `db:"browser" json:"browser"`
	Ip      string `db:"ip" json:"ip"`
}

// 发送注册邮件
type MessageSignupEmail struct {
	Email string          `db:"email" json:"email"`
	Type  string          `db:"type" json:"type"`
	Data  json.RawMessage `db:"data" json:"data"`
}
