package store

import (
	"fmt"
	"log"

	rds "github.com/garyburd/redigo/redis"
	"github.com/hiwjd/horn/mysql"
	"github.com/hiwjd/horn/redis"
)

type DefaultStore struct {
	redisManager *redis.Manager
	mysqlManager *mysql.Manager
}

func NewDefaultStore(redisManager *redis.Manager, mysqlManager *mysql.Manager) Store {
	return &DefaultStore{
		redisManager: redisManager,
		mysqlManager: mysqlManager,
	}
}

func (s *DefaultStore) GetUidsByChatId(chatId string) []string {
	db, err := s.mysqlManager.Get("read")
	if err != nil {
		log.Printf(" 获取mysql连接失败: %s \r\n", err.Error())
		return nil
	}

	var rows []string
	err = db.Select(&rows, "select uid from chat_users where chat_id = ?", chatId)
	if err != nil {
		log.Printf(" mysql执行失败: %s \r\n", err.Error())
		return nil
	}

	return rows
}

func (s *DefaultStore) GetPushAddrByUid(uid string) string {
	conn, err := s.redisManager.Get("node1")
	if err != nil {
		log.Printf(" 获取redis连接失败: %s \r\n", err.Error())
		return ""
	}
	defer conn.Close()

	addr, err := rds.String(conn.Do("GET", fmt.Sprintf("uid-pusher-addr-%s", uid)))
	if err != nil {
		log.Printf(" redis执行失败: %s \r\n", err.Error())
		return ""
	}

	return addr
}
