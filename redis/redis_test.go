package redis

import (
	"testing"

	"github.com/garyburd/redigo/redis"
)

var (
	cfgs    map[string]*Config
	manager *Manager
)

func init() {
	cfgs = make(map[string]*Config, 2)
	cfgs["node1"] = &Config{
		Addr: "127.0.0.1:6379",
		Pass: "",
	}
	cfgs["node2"] = &Config{
		Addr: "127.0.0.1:6379",
		Pass: "",
	}

	manager = New(cfgs)
}

func TestConfig(t *testing.T) {
	if manager.cfgs["node1"] != cfgs["node1"] || manager.cfgs["node2"] != cfgs["node2"] {
		t.Error("数据库配置设置不正确")
	}

	cfg, err := manager.getCfg("node1")
	if err != nil {
		t.Error("应该可以获取到配置")
	}
	if manager.cfgs["node1"] != cfg {
		t.Error("获取到的配置不对")
	}

	cfg, err = manager.getCfg("node2")
	if err != nil {
		t.Error("应该可以获取到配置")
	}
	if manager.cfgs["node2"] != cfg {
		t.Error("获取到的配置不对")
	}
}

func TestMain(t *testing.T) {
	conn, err := manager.Get("node1")
	if err != nil {
		t.Error(err)
	}

	reply, err := redis.String(conn.Do("ping"))
	if err != nil {
		t.Error(err)
	}

	if reply != "PONG" {
		t.Error("PING的返回结果应该是PONG")
	}
}
