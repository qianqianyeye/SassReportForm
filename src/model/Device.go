package model

import (
	"SaasServiceGo/src/webgo"
)

type Device struct {
	ID           int64              `gorm:"column:id" json:"id"`
	HardwareId   string             `gorm:"column:hardware_id" json:"hardwareId"`
	ModelId      int64              `gorm:"column:model_id" json:"modelId"`
	IsBind       int64              `gorm:"column:is_bind" json:"isBind"`
	ReleaseDate  webgo.JsonDateTime `gorm:"column:release_date" json:"releaseDate"`
	StoreId      int64              `gorm:"column:store_id" json:"storeId"`
	Status       int64              `gorm:"column:status" json:"status"`
	CreateTime   webgo.JsonDateTime `gorm:"column:create_time" json:"createTime"`
	UpdateTime   webgo.JsonDateTime `gorm:"column:update_time" json:"updateTime"`
	DeviceNumber string             `gorm:"column:device_number" json:"deviceNumber"`
	Alias        string             `gorm:"column:alias" json:"alias"`
	OnlineDevice string             `json:"onlineDevice"`

	OtherData map[string]interface{} `json:"-"`
	DeviceModular []DeviceModular   `gorm:"ForeignKey:DeviceID;AssociationForeignKey:ID;" json:"device_modular"`
	CommWeimaqi []CommWeimaqi `gorm:"many2many:device_modular;ForeignKey:modular_id;" json:"comm_weimaqi,omitempty"`

	//CommWeimaqi []CommWeimaqi `gorm:"many2many:device_modular;association_jointable_foreignkey:modular_id;" json:"comm_weimaqi,omitempty"`
	CommLgz []CommLgz `gorm:"many2many:device_modular;association_jointable_foreignkey:modular_id;" json:"comm_lgz,omitempty"`
}

func (Device) TableName() string {
	return "device"
}
