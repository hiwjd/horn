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
	Id        string    `json:"sid" db:"sid"`
	Cid       int       `json:"oid" db:"oid"`
	Name      int       `json:"name" db:"name"`
	Gender    int       `json:"gender" db:"gender"`
	Mobile    int       `json:"mobile" db:"mobile"`
	Email     int       `json:"email" db:"email"`
	Pass      int       `json:"pass" db:"pass"`
	Tel       int       `json:"tel" db:"tel"`
	QQ        int       `json:"qq" db:"qq"`
	Status    int       `json:"status" db:"status"`
	State     int       `json:"state" db:"state"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Visitor struct {
	Id        string    `json:"vid" db:"vid"`
	Cid       string    `json:"oid" db:"oid"`
	State     string    `json:"state" db:"state"`
	Fp        string    `json:"fp" db:"fp"`
	tid       string    `json:"tid" db:"tid"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type State interface {
	StaffOnline(oid int, mid string, sid string) error
	StaffOffline(oid int, mid string, sid string) error
	VisitorOnline(oid int, mid string, vid string) error
	VisitorOffline(oid int, mid string, vid string) error
	CreateChat(oid int, mid string, cid, uid string) error
	JoinChat(oid int, mid string, cid, uid string) error
	LeaveChat(oid int, mid string, cid, uid string) error
	GetUidsInChat(oid int, cid string) ([]string, error)
	OnlineStaffList(oid int) ([]*Staff, error)
	GetChatIdsByUid(oid int, uid string) ([]string, error)
	GetPushAddrByUid(oid int, uid string) (string, error)
	GetSidsInOrg(oid int) ([]string, error)
}
