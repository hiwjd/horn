package redis

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	ErrConfigNotFound = errors.New("config not found")
)

type Manager struct {
	cfgs  map[string]*Config
	pools map[string]*redis.Pool
	mutex sync.Mutex
}

func New(cfgs map[string]*Config) *Manager {
	l := len(cfgs)
	pools := make(map[string]*redis.Pool, l)
	var mutex sync.Mutex
	return &Manager{cfgs: cfgs, pools: pools, mutex: mutex}
}

func (m *Manager) Get(tag string) (redis.Conn, error) {
	if pool, ok := m.pools[tag]; ok {
		return pool.Get(), nil
	}

	cfg, err := m.getCfg(tag)
	if err != nil {
		return nil, err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	pool := NewPool(cfg)
	m.pools[tag] = pool

	return pool.Get(), nil
}

func (m *Manager) getCfg(tag string) (*Config, error) {
	if cfg, ok := m.cfgs[tag]; ok {
		return cfg, nil
	}

	return nil, ErrConfigNotFound
}

func NewPool(cfg *Config) *redis.Pool {
	return &redis.Pool{
		MaxActive:   3,
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", cfg.Addr)
			if err != nil {
				return nil, err
			}

			if cfg.Pass != "" {
				if _, err := c.Do("AUTH", cfg.Pass); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

// keyPusher 返回某个Pusher对应在redis里的键
func keyPusher(id string) string {
	return fmt.Sprintf("pusher-%s", id)
}

// keySortedPusher 返回排序的Pusher列表键
func keySortedPusher() string {
	return "sorted-pusher"
}

// keyUser 返回某个用户对应在redis里的键
func keyUser(id string) string {
	return fmt.Sprintf("user-%s", id)
}

// 对话
func keyChat(id string) string {
	return fmt.Sprintf("chat-%s", id)
}

// 对话参与的用户
func keyChatUsers(id string) string {
	return fmt.Sprintf("chat-user-%s", id)
}

// 用户参与的对话
func keyUserChats(id string) string {
	return fmt.Sprintf("user-chats-%s", id)
}

// 指纹对应的uid
func keyFPUID(fp string) string {
	return fmt.Sprintf("user-fp-%s", fp)
}
