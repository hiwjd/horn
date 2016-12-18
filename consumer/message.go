package consumer

import (
	"encoding/json"

	"time"

	"github.com/hiwjd/horn/state"
)

// 消息来源方信息
type From struct {
	Oid  int    `db:"oid" json:"oid"`
	Uid  string `db:"uid" json:"uid"`
	Name string `db:"name" json:"name"`
	Role string `db:"role" json:"role"`
}

// 对话信息
type Chat struct {
	Cid     string         `db:"cid" json:"cid"`
	Oid     int            `db:"oid" json:"oid"`
	Vid     string         `db:"vid" json:"vid"`
	Sid     string         `db:"sid" json:"sid"`
	Tid     string         `db:"tid" json:"tid"`
	Visitor *state.Visitor `db:"visitor" json:"visitor"`
	Staff   *state.Staff   `db:"staff" json:"staff"`
	Tracks  []*state.Track `db:"tracks" json:"tracks"`
}

// // 访客
// type Visitor struct {
// 	Oid  int    `db:"oid" json:"oid"`
// 	Vid  string `db:"vid" json:"vid"`
// 	Name string `db:"name" json:"name"`
// 	Fp   string `db:"fp" json:"fp"`
// }

// // 客服
// type Staff struct {
// 	Oid  int    `db:"oid" json:"oid"`
// 	Sid  string `db:"sid" json:"sid"`
// 	Name string `db:"name" json:"name"`
// }

// 消息基本信息 会匿名组合到具体的消息里
type Message struct {
	Type      string         `db:"type" json:"type"`             // 消息的类型 text, file, image, event
	T         map[string]int `db:"t" json:"t"`                   // 0:客户端发出时间戳 1:入队列时间戳 2:分发时间戳
	Mid       string         `db:"mid" json:"mid"`               // 消息ID
	From      *From          `db:"from" json:"from"`             // 消息发送方信息
	Oid       int            `db:"oid" json:"oid"`               // 公司ID
	CreatedAt time.Time      `db:"created_at" json:"created_at"` // 创建时间戳
}

// 普通消息
type MessageText struct {
	Message
	Cid  string `db:"cid" json:"cid"` // 对话ID
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
	Cid   string `db:"cid" json:"cid"` // 对话ID
	Image *Image `db:"image" json:"image"`
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
	Cid  string `db:"cid" json:"cid"` // 对话ID
	File *File  `db:"file" json:"file"`
}

// 请求对话数据
type EventRequestChat struct {
	Chat *Chat `db:"chat" json:"chat"` // 对话信息
}

// 请求对话消息
type MessageEventRequestChat struct {
	Message
	Event *EventRequestChat `db:"event" json:"event"`
}

// 加入对话数据
type EventJoinChat struct {
	Cid string `db:"cid" json:"cid"` // 对话ID
}

// 加入对话消息
type MessageEventJoinChat struct {
	Message
	Event *EventJoinChat `db:"event" json:"event"`
}

type Message2Pusher struct {
	Type string      `db:"type" json:"type"`
	To   []string    `db:"to" json:"to"`
	Data interface{} `db:"data" json:"data"`
}

// // 访问数据
// type Track struct {
// 	Tid     string `db:"tid" json:"tid"`
// 	Vid     string `db:"vid" json:"vid"`
// 	Fp      string `db:"fb" json:"fp"`
// 	Oid     int    `db:"oid" json:"oid"`
// 	Url     string `db:"url" json:"url"`
// 	Title   string `db:"title" json:"title"`
// 	Referer string `db:"referer" json:"referer"`
// 	Os      string `db:"os" json:"os"`
// 	Browser string `db:"browser" json:"browser"`
// 	Ip      string `db:"ip" json:"ip"`
// 	Addr    string `db:"addr" josn:"addr"`
// }

// 发送注册邮件
type MessageSendEmail struct {
	Email string          `db:"email" json:"email"`
	Type  string          `db:"type" json:"type"`
	Data  json.RawMessage `db:"data" json:"data"`
}

// 超时消息
type MessageTimeout struct {
	Type string `db:"type" json:"type"`
	Oid  int    `db:"oid" json:"oid"`
	Uid  string `db:"uid" json:"uid"`
	Mid  string `db:"mid" json:"mid"`
}
