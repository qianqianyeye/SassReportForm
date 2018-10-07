package model

import (
	"SaasServiceGo/src/webgo"
)

type CommWeimaqi struct {
	ID              int64              `gorm:"column:id" json:"id"`
	WeimaqiId       string             `gorm:"column:weimaqi_id" json:"weimaqiId"`
	Tag             string             `gorm:"column:tag" json:"tag"`
	MarkeType       string             `gorm:"column:marke_type" json:"markeType"`
	Rssi            string             `gorm:"column:rssi" json:"rssi"`
	IsBind          int64              `gorm:"column:is_bind" json:"isBind"`
	HardwareVersion string             `gorm:"column:hardware_version" json:"hardwareVersion"`
	NetworkingType  int64              `gorm:"column:networking_type" json:"networkingType"`
	Status          int64              `gorm:"column:status" json:"status"`
	CreateTime      webgo.JsonDateTime `gorm:"column:create_time" json:"createTime"`
	UpdateTime      webgo.JsonDateTime `gorm:"column:update_time" json:"updateTime"`

	//Device []Device `gorm:"many2many:device_modular;"`
}

func (CommWeimaqi) TableName() string {
	return "comm_weimaqi"
}

