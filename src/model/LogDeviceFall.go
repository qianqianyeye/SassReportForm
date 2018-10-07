package model

import "SaasServiceGo/src/webgo"

type LogDeviceFall struct {
	ID            int64              `gorm:"column:id" json:"id"`
	DeviceId      int64              `gorm:"column:device_id" json:"deviceId"`
	StoreId       int64              `gorm:"column:store_id" json:"storeId"`
	MerchantId    int64              `gorm:"column:merchant_id" json:"merchantId"`
	ModularId     int64              `gorm:"column:modular_id" json:"modularId"`
	ModularNumber string             `gorm:"column:modular_number" json:"modularNumber"`
	FallType      int64              `gorm:"column:fall_type" json:"fallType"`
	FallNum       int64              `gorm:"column:fall_num" json:"fallNum"`
	CreateTime    webgo.JsonDateTime `gorm:"column:create_time" json:"createTime"`
	ModularType   int64              `gorm:"column:modular_type" json:"modularType"`

	RelDeviceFallRftags []RelDeviceFallRftags `gorm:"ForeignKey:DevicefallId"`
	LogDeviceFallMember LogDeviceFallMember `gorm:"ForeignKey:LogDeviceFallId;AssociationForeignKey:ID;" json:"log_device_fall_member,omitempty"`
}

//设置默认表名
func (LogDeviceFall) TableName() string {
	return "log_device_fall"
}
