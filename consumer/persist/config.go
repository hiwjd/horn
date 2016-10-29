package persist

import (
	"github.com/hiwjd/horn/mysql"
)

type Config struct {
	Channel          string
	Topic            string
	NsqdTCPAddrs     []string
	LookupdHTTPAddrs []string
	MaxInFlight      int
	MysqlConfigs     map[string]*mysql.Config
}
