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
	tid := ""
	vid := utils.RandString(23)

	err := ss.CreateChat(oid, mid, cid, sid, sid, vid, tid)
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
	tid := ""
	sid := utils.RandString(19)

	err := ss.CreateChat(oid, mid, cid, vid, sid, vid, tid)
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
	tid := ""

	err := ss.CreateChat(oid, mid, cid, sid, sid, vid, tid)
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

func TestGetStaff(t *testing.T) {
	oid := 1
	sid := "3rUUyOImiv0c2JKelNc"

	staff, err := ss.GetStaff(oid, sid)
	assert.Equal(t, nil, err, "获取客服信息应该成功")
	assert.NotNil(t, staff, "获取客服信息应该成功，不是nil")
}

func TestGetStaffNotExists(t *testing.T) {
	oid := 1
	sid := "xxx"

	staff, err := ss.GetStaff(oid, sid)
	assert.NotEqual(t, nil, err, "获取客服信息应该失败")
	assert.Nil(t, staff, "获取客服信息应该失败，是nil")
}

func TestGetVisitor(t *testing.T) {
	oid := 1
	vid := "SFnvhYMhKzIb9sIaVuvCN9H"

	staff, err := ss.GetVisitor(oid, vid)
	assert.Equal(t, nil, err, "获取访客信息应该成功")
	assert.NotNil(t, staff, "获取访客信息应该成功，不是nil")
}

func TestGetVisitorNotExists(t *testing.T) {
	oid := 1
	vid := "xxx"

	staff, err := ss.GetVisitor(oid, vid)
	assert.NotEqual(t, nil, err, "获取访客信息应该失败")
	assert.Nil(t, staff, "获取访客信息应该失败，是nil")
}

// INSERT INTO `tracks` (`tid`, `vid`, `fp`, `oid`, `url`, `title`, `referer`, `os`, `browser`, `ip`, `addr`, `created_at`)
// VALUES
// 	('20161212140801SFnvhYMhKzIb9sIaVuvCN9HpE4Sb5AIkM7tLRrL', 'SFnvhYMhKzIb9sIaVuvCN9H', 'd9e18d4faaca9be1a797d2d9c4528bce', 1, 'http://www.horn.com:9092/demo.html', 'demo.html', '', 'MacOS', 'Chrome', '101.71.255.226', '浙江杭州', '2016-12-12 14:08:01'),
// 	('20161212140810SFnvhYMhKzIb9sIaVuvCN9H477WvNRW0ASwGYJK', 'SFnvhYMhKzIb9sIaVuvCN9H', '99c5c7c54c42a562e880bb0f522d84d1', 1, 'http://www.horn.com:9092/demo.html', 'demo.html', '', 'MacOS', 'Chrome', '101.71.255.226', '浙江杭州', '2016-12-12 14:08:10'),
// 	('20161212140812SFnvhYMhKzIb9sIaVuvCN9HmubgLIvavDX6e6PE', 'SFnvhYMhKzIb9sIaVuvCN9H', '99c5c7c54c42a562e880bb0f522d84d1', 1, 'http://www.horn.com:9092/demo.html?key=1481551691531', 'demo.html', 'http://www.horn.com:9092/demo.html', 'MacOS', 'Chrome', '101.71.255.226', '浙江杭州', '2016-12-12 14:08:12'),
// 	('20161212140904SFnvhYMhKzIb9sIaVuvCN9HBoTJDV57wZO8NzqJ', 'SFnvhYMhKzIb9sIaVuvCN9H', '99c5c7c54c42a562e880bb0f522d84d1', 1, 'http://www.horn.com:9092/demo.html?key=1481551743671', 'demo.html', 'http://www.horn.com:9092/demo.html?key=1481551691531', 'MacOS', 'Chrome', '101.71.255.226', '浙江杭州', '2016-12-12 14:09:04'),
// 	('20161212140907SFnvhYMhKzIb9sIaVuvCN9HnHkEyuuwuZGm29Lq', 'SFnvhYMhKzIb9sIaVuvCN9H', '99c5c7c54c42a562e880bb0f522d84d1', 1, 'http://www.horn.com:9092/demo.html?key=1481551746955', 'demo.html', 'http://www.horn.com:9092/demo.html?key=1481551743671', 'MacOS', 'Chrome', '101.71.255.226', '浙江杭州', '2016-12-12 14:09:07');

func TestGetVisitorTracks(t *testing.T) {
	oid := 1
	vid := "SFnvhYMhKzIb9sIaVuvCN9H"

	tracks, err := ss.GetVisitorLastTracks(oid, vid, 5)
	assert.Equal(t, nil, err, "获取访客轨迹应该成功")
	assert.NotNil(t, tracks, "获取访客轨迹应该成功，不是nil")
	assert.Equal(t, 5, len(tracks), "访客轨迹应该有5条")
}
