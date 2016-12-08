package state

import (
	"log"

	"github.com/hiwjd/horn/mysql"
	"github.com/hiwjd/horn/redis"
)

type staff struct {
	mysqlManager *mysql.Manager
	redisManager *redis.Manager
}

func (s *staff) online(c *ctx, sid string) error {
	db, err := s.mysqlManager.Get("write")
	if err != nil {
		log.Printf(" 获取mysql连接失败: %s \r\n", err.Error())
		return err
	}

	sql := "update staff set state='on' where oid=? and sid=?"
	r, err := db.Exec(sql, c.oid, sid)
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

func (s *staff) offline(c *ctx, sid string) error {
	db, err := s.mysqlManager.Get("write")
	if err != nil {
		log.Printf(" 获取mysql连接失败: %s \r\n", err.Error())
		return err
	}

	sql := "update staff set state='off' where oid=? and sid=?"
	r, err := db.Exec(sql, c.oid, sid)
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

func (s *staff) onlineStaffList(c *ctx) ([]*Staff, error) {
	db, err := s.mysqlManager.Get("write")
	if err != nil {
		log.Printf(" 获取mysql连接失败: %s \r\n", err.Error())
		return nil, err
	}

	var staffs []*Staff
	sql := "select * from staff where oid = ? and state = 'on'"
	err = db.Select(&staffs, sql, c.oid)
	if err != nil {
		return nil, err
	}

	return staffs, nil
}
