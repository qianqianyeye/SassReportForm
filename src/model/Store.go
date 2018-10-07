package model

import (
	"SaasServiceGo/src/webgo"
)

type Store struct {
	ID              int64              `gorm:"column:id" json:"id"`
	Name            string             `gorm:"column:name" json:"name"`
	AddressProvince string             `gorm:"column:address_province" json:"addressProvince"`
	AddressCity     string             `gorm:"column:address_city" json:"addressCity"`
	AddressArea     string             `gorm:"column:address_area" json:"addressArea"`
	AddressInfo     string             `gorm:"column:address_info" json:"addressInfo"`
	Remarks         string             `gorm:"column:remarks" json:"remarks"`
	Status          int64              `gorm:"column:status" json:"status"`
	CreateTime      webgo.JsonDateTime `gorm:"column:create_time" json:"createTime"`
	UpdateTime      webgo.JsonDateTime `gorm:"column:update_time" json:"updateTime"`
	IsBind          int64              `gorm:"column:is_bind" json:"isBind"`
	ContactId       int64              `gorm:"column:contact_id" json:"contactId"`
}

func (Store) TableName() string {
	return "store"
}
