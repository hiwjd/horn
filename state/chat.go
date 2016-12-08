package state

import (
	"fmt"
	"log"

	rds "github.com/garyburd/redigo/redis"
	"github.com/hiwjd/horn/mysql"
	"github.com/hiwjd/horn/redis"
)

type chat struct {
	mysqlManager *mysql.Manager
	redisManager *redis.Manager
}

func (s *chat) create(c *ctx, cid, uid string) error {
	db, err := s.mysqlManager.Get("write")
	if err != nil {
		log.Printf(" 获取mysql连接失败: %s \r\n", err.Error())
		return err
	}

	var vid, sid string
	if isSid(uid) {
		sid = uid
	} else if isVid(uid) {
		vid = uid
	} else {
		return ErrInvalidUid
	}

	sql := `
	INSERT INTO
		chats(cid,oid,creator,vid,sid,user_num,state)
	VALUES
		(?,?,?,?,?,1,'request')
	`
	r, err := db.Exec(sql, cid, c.oid, uid, vid, sid)
	if err != nil {
		log.Printf(" 创建对话失败: %s \r\n", err.Error())
		return err
	}

	n, err := r.RowsAffected()
	if err != nil {
		return err
	}

	if n != 1 {
		return ErrUpdateNoAffect
	}

	sql = `
	INSERT INTO 
		chat_user(cid, uid, role)
	VALUES 
		(?,?,?)
	`
	role := getRole(uid)
	r, err = db.Exec(sql, cid, uid, role)
	if err != nil {
		log.Printf(" 建立对话-用户关系失败[%s - %s - %s]: %s \r\n", cid, uid, role, err.Error())
		return err
	}

	n, err = r.RowsAffected()
	if err != nil {
		return err
	}

	if n != 1 {
		return ErrUpdateNoAffect
	}

	return nil
}

func (s *chat) addUser(c *ctx, cid, uid string) error {
	db, err := s.mysqlManager.Get("write")
	if err != nil {
		log.Printf(" 获取mysql连接失败: %s \r\n", err.Error())
		return err
	}

	sql := `
	INSERT INTO 
		chat_user(cid, uid, role)
	VALUES 
		(?,?,?)
	ON DUPLICATE KEY UPDATE state = 'join'
	`
	role := getRole(uid)
	_, err = db.Exec(sql, cid, uid, role)
	if err != nil {
		log.Printf(" 建立对话-用户关系失败[%s - %s - %s]: %s \r\n", cid, uid, role, err.Error())
		return err
	}

	if role == "staff" {
		sql = `UPDATE chats SET sid=? WHERE cid = ? and sid = ''`
	} else {
		sql = `UPDATE chats SET vid=? WHERE cid = ? and vid = ''`
	}
	_, err = db.Exec(sql, uid, cid)
	if err != nil {
		log.Printf(" 更新对话用户 角色[%s] ID[%s] 失败: %s \r\n", role, uid, err.Error())
		return err
	}

	sql = `UPDATE chats SET user_num=user_num+1,state='active' WHERE cid = ?`
	_, err = db.Exec(sql, cid)
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

	reply, err := rds.String(conn.Do("set", fmt.Sprintf("state-version-%s", uid), c.mid))
	if err != nil {
		log.Printf(" 执行redis失败: %s \r\n", err.Error())
		return err
	}

	log.Printf(" 执行结果: %s \r\n", reply)

	return nil
}

func (s *chat) removeUser(c *ctx, cid, uid string) error {
	db, err := s.mysqlManager.Get("write")
	if err != nil {
		log.Printf(" 获取mysql连接失败: %s \r\n", err.Error())
		return nil
	}

	sql := "update chat_user set state = 'leave' where cid = ? and uid = ?"
	r, err := db.Exec(sql, cid, uid)
	if err != nil {
		return err
	}

	n, err := r.RowsAffected()
	if err != nil {
		return err
	}

	if n != 1 {
		return ErrUpdateNoAffect
	}

	return nil
}

func (s *chat) getUidsInChat(c *ctx, cid string) ([]string, error) {
	db, err := s.mysqlManager.Get("read")
	if err != nil {
		log.Printf(" 获取mysql连接失败: %s \r\n", err.Error())
		return nil, err
	}

	var rows []string
	err = db.Select(&rows, "select uid from chat_user where cid = ? and state = 'join'", cid)
	if err != nil {
		log.Printf(" mysql执行失败: %s \r\n", err.Error())
		return nil, err
	}

	return rows, nil
}

func (s *chat) getChatIdsByUid(c *ctx, uid string) ([]string, error) {
	db, err := s.mysqlManager.Get("read")
	if err != nil {
		log.Printf(" 获取mysql连接失败: %s \r\n", err.Error())
		return nil, err
	}

	var rows []string
	err = db.Select(&rows, "select cid from chat_user where uid = ? and state = 'join'", uid)
	if err != nil {
		log.Printf(" mysql执行失败: %s \r\n", err.Error())
		return nil, err
	}

	return rows, nil
}

func (s *chat) getPushAddrByUid(c *ctx, uid string) (string, error) {
	conn, err := s.redisManager.Get("node1")
	if err != nil {
		log.Printf(" 获取redis连接失败: %s \r\n", err.Error())
		return "", err
	}
	defer conn.Close()

	addr, err := rds.String(conn.Do("GET", fmt.Sprintf("uid-pusher-addr-%s", uid)))
	if err != nil {
		log.Printf(" redis执行失败: %s \r\n", err.Error())
		return "", err
	}

	return addr, nil
}
