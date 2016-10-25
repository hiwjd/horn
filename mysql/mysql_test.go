package mysql

import "testing"

var (
	cfgs    map[string]*Config
	manager *Manager
)

func init() {
	cfgs = make(map[string]*Config, 2)
	cfgs["write"] = &Config{
		User:     "horn",
		Pass:     "HornMima1@#",
		Addr:     "",
		Protocol: "",
		Dbname:   "",
	}
	cfgs["read"] = &Config{
		User:     "horn",
		Pass:     "HornMima1@#",
		Addr:     "",
		Protocol: "",
		Dbname:   "",
	}

	manager = New(cfgs)
}

func TestConfig(t *testing.T) {
	if manager.cfgs["write"] != cfgs["write"] || manager.cfgs["read"] != cfgs["read"] {
		t.Error("数据库配置设置不正确")
	}

	cfg, err := manager.getCfg("write")
	if err != nil {
		t.Error("应该可以获取到配置")
	}
	if manager.cfgs["write"] != cfg {
		t.Error("获取到的配置不对")
	}

	cfg, err = manager.getCfg("read")
	if err != nil {
		t.Error("应该可以获取到配置")
	}
	if manager.cfgs["read"] != cfg {
		t.Error("获取到的配置不对")
	}
}

func TestMain(t *testing.T) {
	db, err := manager.Get("write")
	if err != nil {
		t.Error(err)
	}

	err = db.Ping()
	if err != nil {
		t.Error(err)
	}
}
