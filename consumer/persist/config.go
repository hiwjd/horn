package persist

import (
	"github.com/hiwjd/horn/mysql"
)

type Config struct {
	Channel          string
	Topics           []string
	NsqdTCPAddrs     []string
	LookupdHTTPAddrs []string
	MaxInFlight      int
	MysqlConfigs     map[string]*mysql.Config
	SendCloudApiUser string
	SendCloudApiKey  string
	SignupTpl        string
	ResetpassTpl     string
}
