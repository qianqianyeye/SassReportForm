package model

import (
	"SaasServiceGo/src/webgo"
)

type CoinLog struct {
	ID            int64              `gorm:"column:id" json:"id"`
	CoinCount     int64              `gorm:"column:coin_count" json:"coinCount"`
	MemberId      int64              `gorm:"column:member_id" json:"memberId"`
	CreateTime    webgo.JsonDateTime `gorm:"column:create_time" json:"createTime"`
	DeviceId      int64              `gorm:"column:device_id" json:"deviceId"`
	StoreId       int64              `gorm:"column:store_id" json:"storeId"`
	MerchantId    int64              `gorm:"column:merchant_id" json:"merchantId"`
	ModularId     int64              `gorm:"column:modular_id" json:"modularId"`
	ModularType   int64              `gorm:"column:modular_type" json:"modularType"`
	Type          int64              `gorm:"column:type" json:"type"`
	Coin          int64              `gorm:"column:coin" json:"coin"`
	Status        int64              `gorm:"column:status" json:"status"`
	OrderId       int64              `gorm:"column:order_id" json:"orderId"`
	OrderType     int64              `gorm:"column:order_type" json:"orderType"`
	InstallSite   int64              `gorm:"column:install_site" json:"installSite"`
	ChooseCount   int64              `gorm:"column:choose_count" json:"chooseCount"`
	OrderNumber   int64              `gorm:"column:order_number" json:"orderNumber"`
	UpdateTime    webgo.JsonDateTime `gorm:"column:update_time" json:"updateTime"`
	DeviceAlias   string             `gorm:"column:device_alias" json:"deviceAlias"`
	Balance       int64              `gorm:"column:balance" json:"balance"`
	MCoinCount    int64              `gorm:"column:m_coin_count" json:"mCoinCount"`
	MBalance      int64              `gorm:"column:m_balance" json:"mBalance"`
	RealCoinCount int64              `gorm:"column:real_coin_count" json:"realCoinCount"`
	Count         int64              `json:"count"`
}

func (CoinLog) TableName() string {
	return "coin_log"
}
