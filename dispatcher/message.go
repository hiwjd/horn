package dispatcher

type From struct {
	Id   string
	Name string
}

type Chat struct {
	Id string
}

type Message struct {
	T    map[string]int // 0:客户端发出时间戳 1:入队列时间戳 2:分发时间戳
	Mid  string
	From From
	Chat Chat
}

type MessageText struct {
	Message
	Text string
}

type Image struct {
	Src    string
	Width  int
	Height int
	Size   int
}

type MessageImage struct {
	Message
	Image Image
}

type File struct {
	Src  string
	Name string
	Size int
}

type MessageFile struct {
	Message
	File File
}

type Message2Pusher struct {
	Type string
	To   []string
	Data interface{}
}
