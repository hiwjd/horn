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
	conn, err := s.redisManager.Get("node1")
	if err != nil {
		log.Printf(" 获取redis连接失败: %s \r\n", err.Error())
		return nil
	}
	defer conn.Close()

	rows, err := rds.Strings(conn.Do("SMEMBERS", fmt.Sprintf("chat-users-%s", chatId)))
	if err != nil {
		log.Printf(" redis执行失败: %s \r\n", err.Error())
		return nil
	}

	return rows

	// db, err := s.mysqlManager.Get("read")
	// if err != nil {
	// 	log.Printf(" 获取mysql连接失败: %s \r\n", err.Error())
	// 	return nil
	// }

	// var rows []string
	// err = db.Select(&rows, "select uid from chat_users where chat_id = ?", chatId)
	// if err != nil {
	// 	log.Printf(" mysql执行失败: %s \r\n", err.Error())
	// 	return nil
	// }

	// return rows
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

func (s *DefaultStore) JoinChat(mid string, chatId string, uid string) error {
	conn, err := s.redisManager.Get("node1")
	if err != nil {
		log.Printf(" 获取redis连接失败: %s \r\n", err.Error())
		return err
	}
	defer conn.Close()

	err = conn.Send("MULTI")
	if err != nil {
		log.Printf(" 执行redis事物失败: %s \r\n", err.Error())
		conn.Do("DISCARD")
		return err
	}
	err = conn.Send("sadd", fmt.Sprintf("user-chats-%s", uid), chatId)
	if err != nil {
		log.Printf(" 执行redis事物失败: %s \r\n", err.Error())
		conn.Do("DISCARD")
		return err
	}
	err = conn.Send("sadd", fmt.Sprintf("chat-users-%s", chatId), uid)
	if err != nil {
		log.Printf(" 执行redis事物失败: %s \r\n", err.Error())
		conn.Do("DISCARD")
		return err
	}
	err = conn.Send("set", fmt.Sprintf("event-version-%s", uid), mid)
	if err != nil {
		log.Printf(" 执行redis事物失败: %s \r\n", err.Error())
		conn.Do("DISCARD")
		return err
	}
	r, err := conn.Do("EXEC")
	if err != nil {
		log.Printf(" 执行redis事物失败: %s \r\n", err.Error())
		conn.Do("DISCARD")
		return err
	}
	log.Printf(" 执行结果: %v \r\n", r)

	return nil
}
