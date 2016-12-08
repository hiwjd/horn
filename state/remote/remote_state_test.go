package remote

import (
	"os"
	"testing"

	"github.com/hiwjd/horn/state"
	"github.com/hiwjd/horn/utils"
	"github.com/stretchr/testify/assert"
)

var (
	c   state.State
	oid int    = 1
	mid string = "fake_mid"
	sid string = "3rUUyOImiv0c2JKelNc"
	vid string = "SFnvhYMhKzIb9sIaVuvCN9H"
)

func TestMain(m *testing.M) {
	c = New("http://127.0.0.1:9094")
	os.Exit(m.Run())
}

func TestStaffOnline(t *testing.T) {
	err := c.StaffOnline(oid, mid, sid)
	assert.Equal(t, nil, err, "客服上线应该成功")
}

func TestStaffOffline(t *testing.T) {
	err := c.StaffOffline(oid, mid, sid)
	assert.Equal(t, nil, err, "客服下线应该成功")
}

func TestVisitorOnline(t *testing.T) {
	err := c.VisitorOnline(oid, mid, vid)
	assert.Equal(t, nil, err, "访客上线应该成功")
}

func TestVisitorOffline(t *testing.T) {
	err := c.VisitorOffline(oid, mid, vid)
	assert.Equal(t, nil, err, "访客下线应该成功")
}

func TestStaffCreateChat_LeaveChat(t *testing.T) {
	cid := utils.RandString(25)
	err := c.CreateChat(oid, mid, cid, sid)
	assert.Equal(t, nil, err, "客服创建对话应该成功")

	cids, err := c.GetChatIdsByUid(oid, sid)
	assert.Equal(t, nil, err, "客服["+sid+"]应该有对话")
	assert.Equal(t, []string{cid}, cids, "客服["+sid+"]应该有对话")

	uids, err := c.GetUidsInChat(oid, cid)
	assert.Equal(t, nil, err, "客服应该在对话参与人中")
	assert.Equal(t, []string{sid}, uids, "客服应该在对话参与人中")

	err = c.LeaveChat(oid, mid, cid, sid)
	assert.Equal(t, nil, err, "客服离开对话应该成功")

	uids, err = c.GetUidsInChat(oid, cid)
	assert.Equal(t, nil, err, "客服不应该在对话参与人中")
	assert.Empty(t, uids, "客服不应该在对话参与人中")

	cids, err = c.GetChatIdsByUid(oid, sid)
	assert.Equal(t, nil, err, "客服["+sid+"]应该没有对话")
	assert.Empty(t, cids, "客服["+sid+"]应该没有对话")
}

func TestVisitorCreateChat_LeaveChat(t *testing.T) {
	cid := utils.RandString(25)
	err := c.CreateChat(oid, mid, cid, vid)
	assert.Equal(t, nil, err, "访客创建对话应该成功")

	cids, err := c.GetChatIdsByUid(oid, vid)
	assert.Equal(t, nil, err, "访客["+vid+"]应该有对话")
	assert.Equal(t, []string{cid}, cids, "访客["+vid+"]应该有对话")

	uids, err := c.GetUidsInChat(oid, cid)
	assert.Equal(t, nil, err, "访客应该在对话参与人中")
	assert.Equal(t, []string{vid}, uids, "访客应该在对话参与人中")

	err = c.LeaveChat(oid, mid, cid, vid)
	assert.Equal(t, nil, err, "访客离开对话应该成功")

	uids, err = c.GetUidsInChat(oid, cid)
	assert.Equal(t, nil, err, "访客不应该在对话参与人中")
	assert.Empty(t, uids, "访客不应该在对话参与人中")

	cids, err = c.GetChatIdsByUid(oid, vid)
	assert.Equal(t, nil, err, "访客["+vid+"]应该没有对话")
	assert.Empty(t, cids, "访客["+vid+"]应该没有对话")
}

func TestGetPushAddrByUid(t *testing.T) {
	addr, err := c.GetPushAddrByUid(oid, vid)
	assert.Equal(t, nil, err, "访客创建对话应该成功")
	assert.NotEmpty(t, addr, "访客创建对话应该成功")
}
