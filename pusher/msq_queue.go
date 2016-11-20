package pusher

import (
	"bytes"
	"errors"
	"log"
	"sync"
)

var (
	// ErrNoNewMsg 代表没有新消息了
	ErrNoNewMsg = errors.New("no new msg")
	// ErrNoMsg 代表没有消息
	ErrNoMsg = errors.New("no msg")
)

// Node 表示单个节点
type Node struct {
	data  []byte // 节点的实际数据
	index int    // 该节点在链表的编号 按入链表的顺序排
	//flag byte // 1:对话消息 2:事件消息
	next *Node // 下一个节点
	prev *Node // 前一个节点
}

// Track 是某个追踪信息
type Track struct {
	index  int       // 最后获得过的节点的编号
	notify chan bool // 是否有新的节点 即是否有编号大雨index的节点
}

// LinkList 消息链表
type LinkList struct {
	sync.RWMutex
	head   *Node // 头节点
	tail   *Node // 尾节点
	size   int   // 最大长度
	iidx   map[int]*Node
	tracks map[string]*Track
}

// NewLinkList 创建新的消息链表
func NewLinkList(size int, trackSize int) *LinkList {
	var lock sync.RWMutex
	iidx := make(map[int]*Node, size)
	tracks := make(map[string]*Track, trackSize)
	return &LinkList{lock, nil, nil, size, iidx, tracks}
}

// Push 追加元素到队尾
func (c *LinkList) Push(data []byte) error {
	c.Lock()
	defer c.Unlock()

	n := &Node{data, 0, nil, nil}
	if c.head == nil {
		c.head = n
		c.tail = n
	} else {
		n.index = c.tail.index + 1
		n.prev = c.tail
		c.tail.next = n
		c.tail = n
	}

	if c.tail.index-c.head.index == c.size {
		// 满了 删除队头
		delete(c.iidx, c.head.index)
		c.head.next.prev = nil
		c.head = c.head.next
	}

	// 维护map索引
	c.iidx[n.index] = n

	// 通知一遍有新消息
	c.notify()

	return nil
}

func (c *LinkList) notify() {
	for _, t := range c.tracks {
		if c.tail.index > t.index {
			// 有新消息的
			select {
			case t.notify <- true:
			default: // 防止t.notify处于阻塞的情况下导致整个通知阻塞
			}
		}
	}
}

// GetNotify 返回消息通知channel
func (c *LinkList) GetNotify(trackID string) <-chan bool {
	var t *Track
	var ok bool
	if t, ok = c.tracks[trackID]; !ok {
		notify := make(chan bool, 1)
		t = &Track{-1, notify}
		c.tracks[trackID] = t
		if c.head != nil {
			notify <- true
		}
	}

	return t.notify
}

// Fetch 获取index之后的所有节点 index是上次获取的最后一个节点的索引值
// 对于之前获取过，然后过了很久才再次来获取，导致中间有部分消息已经不在内存里了，
// 可以通过客户端这样显示来缓解
// 消息0
// 消息1
// 消息2
// [...]（这里可以点击，然后把缺了的这段消息获取回来）
// 消息5
// 消息6
func (c *LinkList) Fetch(trackID string) ([]*Node, error) {
	c.Lock()
	defer c.Unlock()

	var t *Track
	var ok bool
	if t, ok = c.tracks[trackID]; !ok {
		t = &Track{-1, make(chan bool, 1)}
		c.tracks[trackID] = t
	}

	log.Println("")
	log.Printf(" -> trackID:%s \r\n", trackID)
	log.Printf(" -> track:%+v \r\n", t)
	log.Printf(" -> link:%+v \r\n\r\n", c)

	if c.tail == nil {
		return nil, ErrNoMsg
	}

	if t.index >= c.tail.index {
		// 没有最新的了
		return nil, ErrNoNewMsg
	}

	ns := make([]*Node, c.tail.index-t.index)
	n := c.iidx[t.index+1]
	i := 0
	for n != nil {
		ns[i] = n
		n = n.next
		i++
	}

	t.index = c.tail.index

	return ns, nil
}

func (c *LinkList) Del(trackID string) int {
	if _, ok := c.tracks[trackID]; ok {
		delete(c.tracks, trackID)
	}

	return len(c.tracks)
}

// Stats 返回统计数据
func (c *LinkList) Stats() map[string]interface{} {
	m := make(map[string]interface{}, 2)
	m["size"] = c.size

	mt := make(map[string]interface{}, len(c.tracks))
	for tid, t := range c.tracks {
		mt[tid] = t.index
	}
	m["tracks"] = mt

	iidx := make([]map[string]interface{}, len(c.iidx))
	i := 0
	for idx, node := range c.iidx {
		iidx[i] = map[string]interface{}{
			"idx":   idx,
			"index": node.index,
			"data":  node.data,
		}
		i++
	}
	m["iidx"] = iidx

	return m
}

// NodesToJSON 把节点数组转成json数组
func NodesToJSON(nodes []*Node) []byte {
	var buffer bytes.Buffer
	l := len(nodes) - 1
	buffer.WriteString("[")
	for i, n := range nodes {
		if n == nil {
			continue
		}
		buffer.Write(n.data)
		if i != l {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("]")

	return buffer.Bytes()
}
