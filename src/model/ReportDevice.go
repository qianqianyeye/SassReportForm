package model

import "time"

type ReportDevice struct {
	ID               int64     `gorm:"column:id" json:"id"`
	Alias            string    `gorm:"column:alias" json:"alias"`
	ChooseCount      int64     `gorm:"column:choose_count" json:"choose_count"`
	DeviceFallCount  int64     `gorm:"column:device_fall_count" json:"device_fall_count"`
	DeviceNumber     string    `gorm:"column:device_number" json:"device_number"`
	//GameOfflineCount int64     `gorm:"column:game_offline_count" json:"game_offline_count"`
	//GameOnlineCount  int64     `gorm:"column:game_online_count" json:"game_online_count"`
	DeviceId         int64     `gorm:"column:device_id" json:"device_id"`
	LgzCount         int64     `gorm:"column:lgz_count" json:"lgz_count"`
	WeimaqiCount     int64     `gorm:"column:weimaqi_count" json:"weimaqi_count"`
	RealCoinCount    int64     `gorm:"column:real_coin_count" json:"real_coin_count"`
	CoinCount        int64     `gorm:"column:coin_count" json:"coin_count"`
	Status           int64     `gorm:"column:status" json:"status"`
	CreateTime       time.Time `gorm:"column:create_time" json:"-"`
	StoreId          int64     `gorm:"column:store_id" json:"store_id"`
}

func (ReportDevice) TableName() string {
	return "report_device"
}
