package pusher

import (
    "errors"
    "sync"
    "time"
)

/***************************************************
+------------------------------------------------+
|                                                |
|  +-------+                                     |
|  | UID   |                                     |
|  +-------+     track1, offset 1                |
|  |-------| <---+                               |
|  +-------+                                     |
|  |       |                                     |
|  |       |                                     |
|  |       |                                     |
|  +-------+                                     |
|  |-------|                                     |
|  |-------|     track0, offset n                |
|  +-------+ <---+                               |
|  |       |                                     |
|  +---+---+                                     |
|      ^                                Pusher1  |
|      |                                         |
+------------------------------------------------+
       |
       |
       |request with uid and track_id
       +
****************************************************/

var (
    // ErrUIDNotExist 说明下发服务里没有这个UID
    ErrUIDNotExist = errors.New("uid not exist")
    // ErrFetchTimeout 说明获取消息的等待超时了
    ErrFetchTimeout = errors.New("fetch time out")
)

// Pusher 管理一个推送服务器节点的消息下发
type Pusher struct {
    sync.RWMutex
    queues         map[string]*LinkList
    queueSize      int
    queueTrackSize int
}

// New 返回*Pusher
func New(queuesLen, queueSize, queueTrackSize int) *Pusher {
    var lock sync.RWMutex
    queues := make(map[string]*LinkList, queuesLen)
    return &Pusher{lock, queues, queueSize, queueTrackSize}
}

// Add 增加一个消息下发队列
func (c *Pusher) Add(uid string) error {
    c.Lock()
    defer c.Unlock()

    q := NewLinkList(c.queueSize, c.queueTrackSize)
    c.queues[uid] = q
    return nil
}

// Push 追加消息
func (c *Pusher) Push(uid string, data []byte) error {
    var q *LinkList
    var ok bool

    if q, ok = c.queues[uid]; !ok {
        return ErrUIDNotExist
    }

    return q.Push(data)
}

// Fetch 获取消息
func (c *Pusher) Fetch(uid string, trackID string, keep time.Duration) ([]byte, error) {
    var q *LinkList
    var ok bool

    if q, ok = c.queues[uid]; !ok {
        return nil, ErrUIDNotExist
    }

    notify := q.GetNotify(trackID)
    select {
    case <-notify:
        ns, err := q.Fetch(trackID)
        if err != nil {
            return nil, err
        }
        return NodesToJSON(ns), nil
    case <-time.After(keep):
        // timeout
        return nil, ErrFetchTimeout
    }
}

func (c *Pusher) Stats(uid string) map[string]interface{} {
    if uid != "" {
        if q, ok := c.queues[uid]; ok {
            return q.Stats()
        }   
    }

    m := make(map[string]interface{}, 3)
    m["queue_size"] = c.queueSize
    m["track_size"] = c.queueTrackSize

    qs := make(map[string]interface{}, len(c.queues))
    for uid, ll := range c.queues {
        qs[uid] = ll.size
    }
    m["queues"] = qs
    return m
}