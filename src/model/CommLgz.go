package model

import (
	"SaasServiceGo/src/webgo"
)

type CommLgz struct {
	ID int64 `gorm:"column:id" json:"id"`
	LgzId string `gorm:"column:lgz_id" json:"lgz_id"`
	IsBind int64 `gorm:"column:is_bind" json:"is_bind"`
	Status int64 `gorm:"column:status" json:"status"`
	ErrCode string `gorm:"column:err_code" json:"err_code"`
	CreateTime webgo.JsonDateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime webgo.JsonDateTime `gorm:"column:update_time" json:"update_time"`
	Time int64 `gorm:"column:time" json:"time"`
	Strong int64 `gorm:"column:strong" json:"strong"`
	Weak int64 `gorm:"column:weak" json:"weak"`
}

//设置默认表名
func (CommLgz) TableName() string {
	return "comm_lgz"
}
