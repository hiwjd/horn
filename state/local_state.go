package state

import (
	"github.com/hiwjd/horn/mysql"
	"github.com/hiwjd/horn/redis"
)

type localState struct {
	sm *staff
	vm *visitor
	cm *chat
	om *org
}

func New(mysqlManager *mysql.Manager, redisManager *redis.Manager) State {
	sm := &staff{mysqlManager, redisManager}
	vm := &visitor{mysqlManager, redisManager}
	cm := &chat{mysqlManager, redisManager}
	om := &org{mysqlManager, redisManager}

	return &localState{
		sm: sm,
		vm: vm,
		cm: cm,
		om: om,
	}
}

func (s *localState) StaffOnline(oid int, mid string, sid string) error {
	c := &ctx{oid, mid}
	return s.sm.online(c, sid)
}

func (s *localState) StaffOffline(oid int, mid string, sid string) error {
	c := &ctx{oid, mid}
	return s.sm.offline(c, sid)
}

func (s *localState) VisitorOnline(oid int, mid string, vid string) error {
	c := &ctx{oid, mid}
	return s.vm.online(c, vid)
}

func (s *localState) VisitorOffline(oid int, mid string, vid string) error {
	c := &ctx{oid, mid}
	return s.vm.offline(c, vid)
}

func (s *localState) CreateChat(oid int, mid string, cid, uid string) error {
	c := &ctx{oid, mid}
	return s.cm.create(c, cid, uid)
}

func (s *localState) JoinChat(oid int, mid string, cid, uid string) error {
	c := &ctx{oid, mid}
	return s.cm.addUser(c, cid, uid)
}

func (s *localState) LeaveChat(oid int, mid string, cid, uid string) error {
	c := &ctx{oid, mid}
	return s.cm.removeUser(c, cid, uid)
}

func (s *localState) GetUidsInChat(oid int, cid string) ([]string, error) {
	c := &ctx{oid, ""}
	return s.cm.getUidsInChat(c, cid)
}

func (s *localState) OnlineStaffList(oid int) ([]*Staff, error) {
	c := &ctx{oid, ""}
	return s.sm.onlineStaffList(c)
}

func (s *localState) GetChatIdsByUid(oid int, uid string) ([]string, error) {
	c := &ctx{oid, ""}
	return s.cm.getChatIdsByUid(c, uid)
}

func (s *localState) GetPushAddrByUid(oid int, uid string) (string, error) {
	c := &ctx{oid, ""}
	return s.cm.getPushAddrByUid(c, uid)
}

func (s *localState) GetSidsInOrg(oid int) ([]string, error) {
	c := &ctx{oid, ""}
	return s.om.getSidsInOrg(c)
}
