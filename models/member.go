package models

import (
	"OnlineBooks/common"
	"github.com/beego/beego/v2/client/orm"
	"time"
)

type Member struct {
	MemberId      int       `orm:"pk;auto" json:"member_id"`
	Account       string    `orm:"size(30);unique" json:"account"`
	Nickname      string    `orm:"size(30);unique" json:"nickname"`
	Password      string    ` json:"-"`
	Description   string    `orm:"size(640)" json:"description"`
	Email         string    `orm:"size(100);unique" json:"email"`
	Phone         string    `orm:"size(20);null;default(null)" json:"phone"`
	Avatar        string    `json:"avatar"`
	Role          int       `orm:"default(1)" json:"role"`
	RoleName      string    `orm:"-" json:"role_name"`
	Status        int       `orm:"default(0)" json:"status"`
	CreateTime    time.Time `orm:"type(datetime);auto_now_add" json:"create_time"`
	CreateAt      int       `json:"create_at"`
	LastLoginTime time.Time `orm:"type(datetime);null" json:"last_login_time"`
}

func (m *Member) TableName() string {
	return TNMembers()
}

func NewMember() *Member {
	return &Member{}
}

func (m *Member) Find(id int) (*Member, error) {
	m.MemberId = id
	if err := orm.NewOrm().Read(m); err != nil {
		return m, err
	}
	m.RoleName = common.Role(m.Role)
	return m, nil
}

func (m *Member) IsAdministrator() bool {
	if m == nil || m.MemberId <= 0 {
		return false
	}
	return m.Role == 0 || m.Role == 1
}