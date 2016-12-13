package state

import (
	"log"

	"database/sql"

	"fmt"

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

func (s *visitor) getVisitor(c *ctx, vid string) (*Visitor, error) {
	db, err := s.mysqlManager.Get("write")
	if err != nil {
		log.Printf(" 获取mysql连接失败: %s \r\n", err.Error())
		return nil, err
	}

	ss := "select * from visitors where oid=? and vid=?"
	var visitor Visitor
	err = db.Get(&visitor, ss, c.oid, vid)
	if err == sql.ErrNoRows {
		return nil, err
	}

	return &visitor, err
}

func (s *visitor) getVisitorLastTracks(c *ctx, vid string, limit int) ([]*Track, error) {
	db, err := s.mysqlManager.Get("write")
	if err != nil {
		log.Printf(" 获取mysql连接失败: %s \r\n", err.Error())
		return nil, err
	}

	ss := fmt.Sprintf("select * from tracks where oid=? and vid=? order by created_at desc limit %d", limit)
	var tracks []*Track
	err = db.Select(&tracks, ss, c.oid, vid)
	if err == sql.ErrNoRows {
		return nil, err
	}

	return tracks, err
}
