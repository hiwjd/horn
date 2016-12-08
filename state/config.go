package state

import (
	"github.com/hiwjd/horn/mysql"
	"github.com/hiwjd/horn/redis"
)

type Config struct {
	MysqlConfigs map[string]*mysql.Config
	RedisConfigs map[string]*redis.Config
}
