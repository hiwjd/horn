package mysql

import (
	"errors"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	ErrConfigNotFound = errors.New("config not found")
)

type Manager struct {
	cfgs  map[string]*Config
	dbs   map[string]*sqlx.DB
	mutex sync.Mutex
}

func New(cfgs map[string]*Config) *Manager {
	l := len(cfgs)
	dbs := make(map[string]*sqlx.DB, l)
	var mutex sync.Mutex
	return &Manager{cfgs: cfgs, dbs: dbs, mutex: mutex}
}

func (m *Manager) Get(tag string) (*sqlx.DB, error) {
	if h, ok := m.dbs[tag]; ok {
		return h, nil
	}

	cfg, err := m.getCfg(tag)
	if err != nil {
		return nil, err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	h, err := open(cfg)
	if err != nil {
		return nil, err
	}
	m.dbs[tag] = h
	return h, nil
}

func (m *Manager) getCfg(tag string) (*Config, error) {
	if c, ok := m.cfgs[tag]; ok {
		return c, nil
	}

	return nil, ErrConfigNotFound
}

func open(cfg *Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%s:%s@%s(%s)/%s?charset=utf8&collation=utf8_general_ci&parseTime=true",
		cfg.User, cfg.Pass, cfg.Protocol, cfg.Addr, cfg.Dbname)
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(512)
	db.SetMaxIdleConns(256)

	return db, nil
}
