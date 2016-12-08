package state

import (
	"os"
	"testing"

	"github.com/hiwjd/horn/mysql"
	"github.com/hiwjd/horn/redis"
	"github.com/hiwjd/horn/utils"
	"github.com/stretchr/testify/assert"
)

var (
	ss State
)

func TestMain(m *testing.M) {
	redisConfig := map[string]*redis.Config{
		"node1": {
			Addr: "127.0.0.1:6379",
			Pass: "",
		},
		"node2": {
			Addr: "127.0.0.1:6379",
			Pass: "",
		},
	}

	mysqlConfig := map[string]*mysql.Config{
		"write": {
			User:     "horn",
			Pass:     "HornMima1@#",
			Addr:     "127.0.0.1:3306",
			Protocol: "tcp",
			Dbname:   "horn",
		},
		"read": {
			User:     "horn",
			Pass:     "HornMima1@#",
			Addr:     "127.0.0.1:3306",
			Protocol: "tcp",
			Dbname:   "horn",
		},
	}

	// INSERT INTO `staff` (`sid`, `oid`, `name`, `gender`, `mobile`, `email`, `pass`, `tel`, `qq`, `status`, `state`, `created_at`, `updated_at`)
	// VALUES
	// ('3rUUyOImiv0c2JKelNc', 1, '小王', '未知', '', 'swordwinter@126.com', '$2y$10$B.o4IA47k7BQGcfP6ntZoe6huORlPaYiuXZuNEqHflfpC/Zo18DTm', '', '', 'active', 'off', '2016-11-19 08:50:01', '2016-11-27 07:20:22');
	// INSERT INTO `visitors` (`vid`, `oid`, `state`, `fp`, `tid`, `created_at`, `updated_at`)
	// VALUES
	// ('SFnvhYMhKzIb9sIaVuvCN9H', 1, 'on', 'baea1cb71707c8e49eab79b0bfd516a4', '20161121003245SFnvhYMhKzIb9sIaVuvCN9H0wXx1Eqg069CF3bY', '2016-11-20 23:56:57', '2016-11-21 00:32:45');

	redisManager := redis.New(redisConfig)
	mysqlManager := mysql.New(mysqlConfig)
	ss = New(mysqlManager, redisManager)

	os.Exit(m.Run())
}

func TestStaffOnlineSidNotExists(t *testing.T) {
	oid := 1
	mid := "fake_mid"
	sid := "fake_sid"
	err := ss.StaffOnline(oid, mid, sid)

	assert.NotEqual(t, nil, err, "客服上线应该失败")
}

func TestStaffOfflineSidNotExists(t *testing.T) {
	oid := 1
	mid := "fake_mid"
	sid := "fake_sid"
	err := ss.StaffOffline(oid, mid, sid)

	assert.NotEqual(t, nil, err, "客服下线应该失败")
}

func TestStaffOnline(t *testing.T) {
	oid := 1
	mid := "fake_mid"
	sid := "3rUUyOImiv0c2JKelNc"
	err := ss.StaffOnline(oid, mid, sid)

	assert.Equal(t, nil, err, "客服上线应该成功")
}

func TestStaffOffline(t *testing.T) {
	oid := 1
	mid := "fake_mid"
	sid := "3rUUyOImiv0c2JKelNc"
	err := ss.StaffOffline(oid, mid, sid)

	assert.Equal(t, nil, err, "客服下线应该成功")
}

func TestVisitorOnline(t *testing.T) {
	oid := 1
	mid := "fake_mid"
	vid := "SFnvhYMhKzIb9sIaVuvCN9H"
	err := ss.VisitorOnline(oid, mid, vid)

	assert.Equal(t, nil, err, "访客上线应该成功")
}

func TestVisitorOffline(t *testing.T) {
	oid := 1
	mid := "fake_mid"
	vid := "SFnvhYMhKzIb9sIaVuvCN9H"
	err := ss.VisitorOffline(oid, mid, vid)

	assert.Equal(t, nil, err, "访客下线应该成功")
}

func TestStaffCreateChat_LeaveChat(t *testing.T) {
	oid := 1
	mid := "fake_mid"
	cid := utils.RandString(25)
	sid := "3rUUyOImiv0c2JKelNc"

	err := ss.CreateChat(oid, mid, cid, sid)
	assert.Equal(t, nil, err, "客服创建对话应该成功")

	err = ss.LeaveChat(oid, mid, cid, sid)
	assert.Equal(t, nil, err, "客服离开对话应该成功")
}

func TestStaffJoinChat_LeaveChat(t *testing.T) {
	oid := 1
	mid := "fake_mid"
	cid := utils.RandString(25)
	sid := "3rUUyOImiv0c2JKelNc"

	err := ss.JoinChat(oid, mid, cid, sid)
	assert.Equal(t, nil, err, "客服加入对话应该成功")

	err = ss.LeaveChat(oid, mid, cid, sid)
	assert.Equal(t, nil, err, "客服离开对话应该成功")
}

func TestVisitorCreateChat_LeaveChat(t *testing.T) {
	oid := 1
	mid := "fake_mid"
	cid := utils.RandString(25)
	vid := "SFnvhYMhKzIb9sIaVuvCN9H"

	err := ss.CreateChat(oid, mid, cid, vid)
	assert.Equal(t, nil, err, "访客创建对话应该成功")

	err = ss.LeaveChat(oid, mid, cid, vid)
	assert.Equal(t, nil, err, "访客离开对话应该成功")
}

func TestVisitorJoinChat_LeaveChat(t *testing.T) {
	oid := 1
	mid := "fake_mid"
	cid := utils.RandString(25)
	vid := "SFnvhYMhKzIb9sIaVuvCN9H"

	err := ss.JoinChat(oid, mid, cid, vid)
	assert.Equal(t, nil, err, "访客加入对话应该成功")

	err = ss.LeaveChat(oid, mid, cid, vid)
	assert.Equal(t, nil, err, "访客离开对话应该成功")
}

func TestStaffCreateChat_VisitorJoin_GetUids_VistorLeave_GetUids_StaffLeave_GetUids(t *testing.T) {
	oid := 1
	mid := "fake_mid"
	cid := utils.RandString(25)
	sid := "3rUUyOImiv0c2JKelNc"
	vid := "SFnvhYMhKzIb9sIaVuvCN9H"

	err := ss.CreateChat(oid, mid, cid, sid)
	assert.Equal(t, nil, err, "客服创建对话应该成功")

	cids, err := ss.GetChatIdsByUid(oid, sid)
	assert.Equal(t, nil, err, "客服["+sid+"]应该有对话")
	assert.Equal(t, []string{cid}, cids, "客服["+sid+"]应该有对话")

	err = ss.JoinChat(oid, mid, cid, vid)
	assert.Equal(t, nil, err, "访客加入对话应该成功")

	cids, err = ss.GetChatIdsByUid(oid, vid)
	assert.Equal(t, nil, err, "访客["+vid+"]应该有对话")
	assert.Equal(t, []string{cid}, cids, "访客["+vid+"]应该有对话")

	uids, err := ss.GetUidsInChat(oid, cid)
	assert.Equal(t, nil, err, "对话["+cid+"]中的uid不正确")
	assert.Equal(t, []string{sid, vid}, uids, "对话["+cid+"]中的uid不正确")

	err = ss.LeaveChat(oid, mid, cid, vid)
	assert.Equal(t, nil, err, "访客离开对话应该成功")

	cids, err = ss.GetChatIdsByUid(oid, vid)
	assert.Equal(t, nil, err, "访客["+vid+"]应该没有对话")
	assert.Empty(t, cids, "访客["+vid+"]应该没有对话")

	uids, err = ss.GetUidsInChat(oid, cid)
	assert.Equal(t, nil, err, "对话["+cid+"]中的uid不正确")
	assert.Equal(t, []string{sid}, uids, "对话["+cid+"]中的uid不正确")

	err = ss.LeaveChat(oid, mid, cid, sid)
	assert.Equal(t, nil, err, "客服离开对话应该成功")

	cids, err = ss.GetChatIdsByUid(oid, sid)
	assert.Equal(t, nil, err, "客服["+sid+"]应该没有对话")
	assert.Empty(t, cids, "客服["+sid+"]应该没有对话")

	uids, err = ss.GetUidsInChat(oid, cid)
	assert.Equal(t, nil, err, "对话["+cid+"]中的uid不正确")
	assert.Empty(t, uids, "对话["+cid+"]中的uid不正确")
}
