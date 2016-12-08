package state

import (
	"log"

	"github.com/hiwjd/horn/mysql"
	"github.com/hiwjd/horn/redis"
)

type visitor struct {
	mysqlManager *mysql.Manager
	redisManager *redis.Manager
}

func (s *visitor) online(c *ctx, vid string) error {
	db, err := s.mysqlManager.Get("write")
	if err != nil {
		log.Printf(" 获取mysql连接失败: %s \r\n", err.Error())
		return err
	}

	sql := "update visitors set state='on' where oid=? and vid=?"
	r, err := db.Exec(sql, c.oid, vid)
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

func (s *visitor) offline(c *ctx, vid string) error {
	db, err := s.mysqlManager.Get("write")
	if err != nil {
		log.Printf(" 获取mysql连接失败: %s \r\n", err.Error())
		return err
	}

	sql := "update visitors set state='off' where oid=? and vid=?"
	r, err := db.Exec(sql, c.oid, vid)
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
