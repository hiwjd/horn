package state

import "time"

type Org struct {
	Id        int       `json:"oid" db:"oid"`
	Code      string    `json:"code" db:"code"`
	Name      string    `json:"name" db:"name"`
	Mobile    string    `json:"mobile" db:"mobile"`
	Email     string    `json:"email" db:"email"`
	Balance   string    `json:"balance" db:"balance"`
	Status    string    `json:"status" db:"statuscode"`
	Level     string    `json:"level" db:"level"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Staff struct {
	Sid       string    `json:"sid" db:"sid"`
	Oid       int       `json:"oid" db:"oid"`
	Name      string    `json:"name" db:"name"`
	Gender    string    `json:"gender" db:"gender"`
	Mobile    string    `json:"mobile" db:"mobile"`
	Email     string    `json:"email" db:"email"`
	Pass      string    `json:"-" db:"pass"`
	Tel       string    `json:"tel" db:"tel"`
	QQ        string    `json:"qq" db:"qq"`
	Status    string    `json:"status" db:"status"`
	State     string    `json:"state" db:"state"`
	Gid       string    `json:"gid" db:"gid"`
	CCN       int       `json:"ccn" db:"ccn"`
	CCNCUR    int       `json:"ccn_cur" db:"ccn_cur"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Visitor struct {
	Id        string    `json:"vid" db:"vid"`
	Cid       string    `json:"oid" db:"oid"`
	State     string    `json:"state" db:"state"`
	Fp        string    `json:"fp" db:"fp"`
	Tid       string    `json:"tid" db:"tid"`
	Name      string    `json:"name" db:"name"`
	Gender    string    `json:"gender" db:"gender"`
	Age       string    `json:"age" db:"age"`
	Mobile    string    `json:"mobile" db:"mobile"`
	Email     string    `json:"email" db:"email"`
	QQ        string    `json:"qq" db:"qq"`
	Addr      string    `json:"addr" db:"addr"`
	Note      string    `json:"note" db:"note"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Track struct {
	Tid       string    `db:"tid" json:"tid"`
	Vid       string    `db:"vid" json:"vid"`
	Fp        string    `db:"fp" json:"fp"`
	Oid       int       `db:"oid" json:"oid"`
	Url       string    `db:"url" json:"url"`
	Title     string    `db:"title" json:"title"`
	Referer   string    `db:"referer" json:"referer"`
	Os        string    `db:"os" json:"os"`
	Browser   string    `db:"browser" json:"browser"`
	Ip        string    `db:"ip" json:"ip"`
	Addr      string    `db:"addr" json:"addr"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type State interface {
	StaffOnline(oid int, mid string, sid string) error
	StaffOffline(oid int, mid string, sid string) error
	VisitorOnline(oid int, mid string, vid string) error
	VisitorOffline(oid int, mid string, vid string) error
	CreateChat(oid int, mid string, cid, creator, sid, vid, tid string) error
	JoinChat(oid int, mid string, cid, uid string) error
	LeaveChat(oid int, mid string, cid, uid string) error
	GetUidsInChat(oid int, cid string) ([]string, error)
	OnlineStaffList(oid int) ([]*Staff, error)
	GetChatIdsByUid(oid int, uid string) ([]string, error)
	GetPushAddrByUid(oid int, uid string) (string, error)
	GetSidsInOrg(oid int) ([]string, error)
	GetVisitor(oid int, vid string) (*Visitor, error)
	GetStaff(oid int, sid string) (*Staff, error)
	GetVisitorLastTracks(oid int, vid string, limit int) ([]*Track, error)
}
