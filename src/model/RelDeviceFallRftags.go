package model

type RelDeviceFallRftags struct {
	ID           int64  `gorm:"column:id" json:"id"`
	DevicefallId int64  `gorm:"column:devicefall_id" json:"devicefallId"`
	Rftag        string `gorm:"column:rfTag" json:"rfTag"`
}

func (RelDeviceFallRftags) TableName() string {
	return "rel_device_fall_rftags"
}
