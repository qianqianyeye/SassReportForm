package model

import (
	"SaasServiceGo/src/webgo"
)

type MemberVisitLog struct {
	ID         int64              `gorm:"column:id" json:"id"`
	MemberId   int64              `gorm:"column:member_id" json:"memberId"`
	StoreId    int64              `gorm:"column:store_id" json:"storeId"`
	MerchantId int64              `gorm:"column:merchant_id" json:"merchantId"`
	DeviceId   int64              `gorm:"column:device_id" json:"deviceId"`
	TimeDay    string             `gorm:"column:time_day" json:"timeDay"`
	VisitTime  webgo.JsonDateTime `gorm:"column:visit_time" json:"visitTime"`
}

func (MemberVisitLog) TableName() string {
	return "member_visit_log"
}
