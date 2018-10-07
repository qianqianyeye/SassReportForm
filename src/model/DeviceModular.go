package model


import "SaasServiceGo/src/webgo"

type DeviceModular struct {
	ID          int64              `gorm:"column:id" json:"id"`
	DeviceId    int64              `gorm:"column:device_id" json:"deviceId"`
	ModularId   int64              `gorm:"column:modular_id" json:"modularId"`
	ModularType int64              `gorm:"column:modular_type" json:"modularType"`
	InstallSite int64              `gorm:"column:install_site" json:"installSite"`
	ConType     int64              `gorm:"column:con_type" json:"conType"`
	CreateTime  webgo.JsonDateTime `gorm:"column:create_time" json:"createTime"`
	Coin        int64              `gorm:"column:coin" json:"coin"`
	CtrlFall    int64              `gorm:"column:ctrl_fall" json:"ctrlFall"`
	DefaultFall int64              `gorm:"column:default_fall" json:"defaultFall"`
	ActualFall  int64              `gorm:"column:actual_fall" json:"actualFall"`

	CommWeimaqi CommWeimaqi 	   `gorm:"ForeignKey:ID;AssociationForeignKey:ModularId;" json:"comm_weimaqi,omitempty"`
	CommLgz []CommLgz   	`gorm:"ForeignKey:ID;AssociationForeignKey:ModularId;" json:"comm_lgz,omitempty"`
}

//设置默认表名
func (DeviceModular) TableName() string {
	return "device_modular"
}
