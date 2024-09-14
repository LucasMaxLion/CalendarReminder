package model

import "time"

type User struct {
	Id          int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"id"`
	Uid         string    `gorm:"column:uid;type:varchar(59)" json:"uid"`
	Name        string    `gorm:"column:name;type:varchar(255)" json:"name"`
	Password    string    `gorm:"column:password;type:varchar(50)" json:"password"`
	CreatedTime time.Time `gorm:"column:created_time;type:datetime" json:"created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time;type:datetime" json:"updated_time"`
}

func (m *User) TableName() string {
	return "user"
}

type Reminder struct {
	Id          int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"id"`
	Uid         string    `gorm:"column:uid;type:varchar(255)" json:"uid"`
	CreatorId   string    `gorm:"column:creator_id;type:varchar(255)" json:"creator_id"`
	Content     string    `gorm:"column:content;type:varchar(255)" json:"content"`
	RemindTime  time.Time `gorm:"column:remind_time;type:datetime" json:"remind_time"`
	Begin       string    `gorm:"column:begin;type:varchar(255)" json:"begin"`
	Active      string    `gorm:"column:active;type:varchar(255)" json:"active"`
	Email       string    `gorm:"column:email;type:varchar(50)" json:"email"`
	CreatedTime time.Time `gorm:"column:created_time;type:datetime" json:"created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time;type:datetime" json:"updated_time"`
}

func (m *Reminder) TableName() string {
	return "reminder"
}

func (m *Reminder) UpdateReminderStatus(begin, active string) error {
	Conn.Table("reminder").Where("uid = ?", m.Uid).Updates(map[string]any{"begin": begin, "active": active})
	return nil
}
