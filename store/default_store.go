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

func (s *DefaultStore) JoinChat(mid string, chatId string, uid string, role string) error {
	db, err := s.mysqlManager.Get("write")
	if err != nil {
		log.Printf(" 获取mysql连接失败: %s \r\n", err.Error())
		return err
	}

	sql := `
	INSERT INTO 
		chat_user(chat_id, uid, role)
	VALUES 
		(?,?,?)
	`
	_, err = db.Exec(sql, chatId, uid, role)
	if err != nil {
		log.Printf(" 建立对话-用户关系失败: %s \r\n", err.Error())
		return err
	}

	sql = `UPDATE chats SET user_num=user_num+1 WHERE chat_id = ?`
	_, err = db.Exec(sql, chatId)
	if err != nil {
		log.Printf(" 更新对话用户数失败: %s \r\n", err.Error())
		return err
	}

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

func (s *DefaultStore) CreateChat(chatId string, gid string, creator string, kfid int) error {
	db, err := s.mysqlManager.Get("write")
	if err != nil {
		log.Printf(" 获取mysql连接失败: %s \r\n", err.Error())
		return err
	}

	sql := `
	INSERT INTO
		chats(chat_id,gid,creator,kfid,user_num,state)
	VALUES
		(?,?,?,?,1,'request')
	`
	_, err = db.Exec(sql, chatId, gid, creator, kfid)
	if err != nil {
		log.Printf(" 创建对话失败: %s \r\n", err.Error())
		return err
	}

	return err
}

func (s *DefaultStore) GetStaffsByCompany(cid string) []string {
	conn, err := s.redisManager.Get("node1")
	if err != nil {
		log.Printf(" 获取redis连接失败: %s \r\n", err.Error())
		return nil
	}
	defer conn.Close()

	rows, err := rds.Strings(conn.Do("SMEMBERS", fmt.Sprintf("company-staffs-%s", cid)))
	if err != nil {
		log.Printf(" redis执行失败: %s \r\n", err.Error())
		return nil
	}

	return rows
}
