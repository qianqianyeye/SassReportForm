package model

import (
	"SaasServiceGo/src/webgo"
)

type LogWeimaqiCoin struct {
	ID         int64              `gorm:"column:id" json:"id"`
	Sn         string             `gorm:"column:sn" json:"sn"`
	CreateTime webgo.JsonDateTime `gorm:"column:create_time" json:"createTime"`
	WeimaqiId  string             `gorm:"column:weimaqi_id" json:"weimaqiId"`
	Coinin     int64              `gorm:"column:coinin" json:"coinin"`
	Payout     int64              `gorm:"column:payout" json:"payout"`
}

//设置默认表名
func (LogWeimaqiCoin) TableName() string {
	return "log_weimaqi_coin"
}
