package model

import (
	"time"
)

type ReportForm struct {
	ID               int64     `gorm:"column:id" json:"id"`
	CreateTime       time.Time `gorm:"column:create_time" json:"create_time"`
	ChooseCount      int64     `gorm:"column:choose_count" json:"choose_count"`
	DeviceFallCount  int64     `gorm:"column:device_fall_count" json:"device_fall_count"`
	GameOrderCount   int64     `gorm:"column:game_order_count" json:"game_order_count"`
	CoinCount        int64     `gorm:"column:coin_count" json:"coin_count"`
	MemberCount      int64     `gorm:"column:member_count" json:"member_count"`
	DayMemberCount   int64     `gorm:"column:day_member_count" json:"day_member_count"`
	VisitCount       int64     `gorm:"column:visit_count" json:"visit_count"`
	LgzCoinCount     int64     `gorm:"column:lgz_coin_count" json:"lgz_coin_count"`
	WeimaqiCoinCount int64     `gorm:"column:weimaqi_coin_count" json:"weimaqi_coin_count"`
	StoreId          int64     `gorm:"column:store_id" json:"store_id"`
	//CoinStatistics string `gorm:"column:coin_statistics" json:"coin_statistics"`
	RealCoinCount int64  `gorm:"column:real_coin_count" json:"real_coin_count"`
	MCoinCount    int64  `gorm:"column:m_coin_count" json:"m_coin_count"`
	CoinSpec      string `gorm:"column:coin_spec"  json:"coin_spec"`
}

func (ReportForm) TableName() string {
	return "report_form"
}

