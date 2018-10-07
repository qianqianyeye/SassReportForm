package model

import (
	"SaasServiceGo/src/webgo"
)

type MerchantStore struct {
	ID         int64              `gorm:"column:id" json:"id"`
	MerchantId int64              `gorm:"column:merchant_id" json:"merchantId"`
	StoreId    int64              `gorm:"column:store_id" json:"storeId"`
	CreateTime webgo.JsonDateTime `gorm:"column:create_time" json:"createTime"`
}

func (MerchantStore) TableName() string {
	return "merchant_store"
}
