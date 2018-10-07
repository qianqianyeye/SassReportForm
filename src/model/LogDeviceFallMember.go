package model

import "time"

type LogDeviceFallMember struct {
	ID int64 `gorm:"column:id" json:"id"`
	MemberId int64 `gorm:"column:member_id" json:"member_id"`
	LogDeviceFallId int64 `gorm:"column:log_device_fall_id" json:"log_device_fall_id"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	Rftag string `gorm:"column:rftag" json:"rftag"`
	GoodsName string `gorm:"column:goods_name" json:"goods_name"`
	GoodsImg string `gorm:"column:goods_img" json:"goods_img"`
	GoodsSize string `gorm:"column:goods_size" json:"goods_size"`
	MerchantId int64 `gorm:"column:merchant_id" json:"merchant_id"`
	StoreId int64 `gorm:"column:store_id" json:"store_id"`
	GoodsId int64 `gorm:"column:goods_id" json:"goods_id"`
	GoodsType int64 `gorm:"column:goods_type" json:"goods_type"`
}

//设置默认表名
func (LogDeviceFallMember) TableName() string {
	return "log_device_fall_member"
}