package model

import (
	"SaasServiceGo/src/webgo"
)

type LogLgzCoin struct {
	ID          int64              `gorm:"column:id" json:"id"`
	LgzId       string             `gorm:"column:lgz_id" json:"lgzId"`
	Coin        int64              `gorm:"column:coin" json:"coin"`
	CoinInstall int64              `gorm:"column:coin_install" json:"coinInstall"`
	DeviceId    int64              `gorm:"column:device_id" json:"deviceId"`
	CreateTime  webgo.JsonDateTime `gorm:"column:create_time" json:"createTime"`
	Device      Device             `gorm:"ForeignKey:DeviceId;AssociationForeignKey:ID"`
}

//设置默认表名
func (LogLgzCoin) TableName() string {
	return "log_lgz_coin"
}
