package state

import (
	"log"

	"github.com/hiwjd/horn/mysql"
	"github.com/hiwjd/horn/redis"
)

type org struct {
	mysqlManager *mysql.Manager
	redisManager *redis.Manager
}

func (s *org) getSidsInOrg(c *ctx) ([]string, error) {
	db, err := s.mysqlManager.Get("write")
	if err != nil {
		log.Printf(" 获取mysql连接失败: %s \r\n", err.Error())
		return nil, err
	}

	var rows []string
	sql := "select sid from staff where oid = ?"
	err = db.Select(&rows, sql, c.oid)
	if err != nil {
		return nil, err
	}

	return rows, nil
}
