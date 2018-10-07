package model

import (
	"SaasServiceGo/src/webgo"
)

type Member struct {
	ID          int64              `gorm:"column:id" json:"id"`
	Name        string             `gorm:"column:name" json:"name"`
	Mobile      string             `gorm:"column:mobile" json:"mobile"`
	Integral    int64              `gorm:"column:integral" json:"integral"`
	Headimgurl  string             `gorm:"column:headimgurl" json:"headimgurl"`
	CoinCount   int64              `gorm:"column:coin_count" json:"coin_count"`
	Recharge    int64              `gorm:"column:recharge" json:"recharge"`
	PrizeAmount int64              `gorm:"column:prize_amount" json:"prize_amount"`
	Balance     int64              `gorm:"column:balance" json:"balance"`
	Status      int64              `gorm:"column:status" json:"status"`
	CreateTime  webgo.JsonDateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime  webgo.JsonDateTime `gorm:"column:update_time" json:"update_time"`
}

func (Member) TableName() string {
	return "member"
}
