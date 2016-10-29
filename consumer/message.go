package consumer

type From struct {
	Id   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type Chat struct {
	Id string `db:"id" json:"id"`
}

type Message struct {
	Type string         `db:"type" json:"type"` // 消息的类型 text, file, image, cmd
	T    map[string]int `db:"t" json:"t"`       // 0:客户端发出时间戳 1:入队列时间戳 2:分发时间戳
	Mid  string         `db:"mid" json:"mid"`   // 消息ID
	From From           `db:"from" json:"from"` // 消息发送方信息
	Chat Chat           `db:"chat" json:"chat"` // 对话信息
}

type MessageText struct {
	Message
	Text string `db:"text" json:"text"`
}

type Image struct {
	Src    string `db:"src" json:"src"`
	Width  int    `db:"width" json:"width"`
	Height int    `db:"height" json:"height"`
	Size   int    `db:"size" json:"size"`
}

type MessageImage struct {
	Message
	Image Image `db:"image" json:"image"`
}

type File struct {
	Src  string `db:"src" json:"src"`
	Name string `db:"name" json:"name"`
	Size int    `db:"size" json:"size"`
}

type MessageFile struct {
	Message
	File File `db:"file" json:"file"`
}

type Message2Pusher struct {
	Type string      `db:"type" json:"type"`
	To   []string    `db:"to" json:"to"`
	Data interface{} `db:"data" json:"data"`
}
