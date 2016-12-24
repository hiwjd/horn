package state

import (
	"log"

	"github.com/hiwjd/horn/mysql"
	"github.com/hiwjd/horn/redis"
	"github.com/jmoiron/sqlx"
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

func (s *localState) CreateChat(oid int, mid string, cid, creator, sid, vid, tid string) error {
	c := &ctx{oid, mid}
	return s.cm.create(c, cid, creator, sid, vid, tid)
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

func (s *localState) GetVisitor(oid int, vid string) (*Visitor, error) {
	c := &ctx{oid, ""}
	return s.vm.getVisitor(c, vid)
}

func (s *localState) GetStaff(oid int, sid string) (*Staff, error) {
	c := &ctx{oid, ""}
	return s.sm.getStaff(c, sid)
}

func (s *localState) GetVisitorLastTracks(oid int, vid string, limit int) ([]*Track, error) {
	c := &ctx{oid, ""}
	return s.vm.getVisitorLastTracks(c, vid, limit)
}

func manageStaffCCNCur(db *sqlx.DB, oid int, sid string) error {
	log.Printf("维护客服的当前对话数 oid[%d] sid[%s] \r\n", oid, sid)
	sql := `
		UPDATE staff SET 
			ccn_cur = IFNULL((SELECT count(1) FROM chat_user WHERE oid=? AND uid=? AND role='staff' AND state='join'),0) 
		WHERE 
			oid = ? AND sid = ?
	`
	_, err := db.Exec(sql, oid, sid, oid, sid)
	if err != nil {
		log.Printf(" 维护客服当前对话数失败 oid[%d] sid[%s] err[%s]\r\n", oid, sid, err.Error())
		return err
	}

	return nil
}
