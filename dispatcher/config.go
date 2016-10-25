package dispatcher

import (
	"github.com/hiwjd/horn/mysql"
	"github.com/hiwjd/horn/redis"
)

type Config struct {
	Topic            string
	Channel          string
	LookupdHTTPAddrs []string
	NsqdTCPAddrs     []string
	MaxInFlight      int
	MysqlConfigs     map[string]*mysql.Config
	RedisConfigs     map[string]*redis.Config
}
