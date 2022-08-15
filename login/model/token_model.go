package model

import "time"

type UserToken struct {
	Id         uint64    `gorm:"column:id" json:"id"`
	Uid        uint64    `gorm:"column:uid" json:"uid"`
	TokenID    string    `gorm:"column:token_id" json:"token_id"`
	Token      string    `gorm:"column:token" json:"token"`
	ExpireTime time.Time `gorm:"column:expire_time" json:"expire_time"`
	CreatedAt  time.Time `gorm:"column:creaed_at" json:"created_at"`
	DeletedAt  time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (ut *UserToken) TableName() string {
	return "user_token"
}
