package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
	"SaasServiceGo/src/model"
	"SaasServiceGo/src/webgo"
	"SaasServiceGo/src/db"
)

type ReportFormService struct {
}

//func (ctrl *ReportFormService) GetRelMerchantStore(parm map[string]interface{}) bool{
//	storeId :=webgo.GetResult(parm["storeId"])
//	merchantId :=webgo.GetResult(parm["merchantId"])
//	var merchantStore model.MerchantStore
//	row:=db.SqlDB.Where("store_id=? and merchant_id= ?", storeId,merchantId).First(&merchantStore)
////	row := db.SqlDB.Table("merchant_store").Select("*").Where("store_id=? and merchant_id= ? ",storeId,merchantId).Scan(&merchantStore)
//	if row.RowsAffected>0 {
//		return true
//	}else {
//		return false
//	}
//}

func (ctrl *ReportFormService) GetReportForm(parm map[string]interface{}, ctx *gin.Context) map[string]interface{} {
	storeId := webgo.GetResult(parm["storeId"])
	startTime := webgo.GetResult(parm["start_time"])
	endTime := webgo.GetResult(parm["end_time"])
	merchantId := webgo.GetResult(ctx.Keys["merchantId"])
	top := 4
	DeviceMap, dataMap := GetDeviceList(storeId, startTime, endTime, merchantId)
	Visit := GetVisitCount(storeId, merchantId, startTime, endTime)
	CoinLogMap := GetCoinLogByStoreId(storeId, merchantId, startTime, endTime)
	all_device := len(DeviceMap)
	online_device := webgo.GetResult(dataMap["online"])
	visit_count := len(Visit)
	game_order_count := webgo.GetResult(dataMap["gameOrderCount"])
	coin_count := webgo.GetResult(CoinLogMap["coinCount"])
	member := GetCoinLogGroupByForMember(storeId, merchantId, startTime, endTime)
	member_count := len(member)
	var upt float64 = 0
	if member_count != 0 {
		s, _ := strconv.Atoi(game_order_count)
		upt = float64(s) / float64(member_count)
	} else {
		upt = 0
	}
	var pct float64 = 0
	if member_count != 0 {
		s, _ := strconv.Atoi(coin_count)
		pct = float64(s) / float64(member_count)
	} else {
		pct = 0
	}
	coin_statistics := GetCoinStatistice(storeId, merchantId, startTime, endTime, top)
	history := GetHistory(storeId, merchantId, startTime, endTime)
	device_fall_count := getDeviceFallCountByStoreId(storeId, merchantId, startTime, endTime)
	choose_count := CoinLogMap["chooseCount"]
	offline_coin_count := dataMap["allOffLineCoin"]
	online_coin_count := dataMap["allOnlineCoin"]

	var result map[string]interface{}
	result = make(map[string]interface{})
	result["all_device"] = all_device
	result["online_device"] = online_device
	result["visit_count"] = visit_count
	result["game_order_count"] = game_order_count
	result["upt"] = upt
	result["pct"] = pct
	result["coin_statistics"] = coin_statistics
	result["history"] = history
	result["device_fall_count"] = device_fall_count
	result["choose_count"] = choose_count
	result["member_count"] = member_count
	result["devices"] = DeviceMap
	result["offline_coin_count"] = offline_coin_count
	result["online_coin_count"] = online_coin_count
	return result
}

func GetCoinLogByStoreId(storeId string, merchantId string, start string, end string) map[string]interface{} {
	var coinLog []model.CoinLog
	var chooseCount int64
	var coinCount int64
	chooseCount = 0
	coinCount = 0
	resultMap := make(map[string]interface{})
	if storeId != "0" && merchantId != "0" {
		db.SqlDB.Table("coin_log").Select("*").Where("status=? and merchant_id=? and store_id =? and create_time BETWEEN ? and ?", "1", merchantId, storeId, start, end).Scan(&coinLog)
	} else if storeId == "0" && merchantId != "0" {
		db.SqlDB.Table("coin_log").Select("*").Where("status=? and merchant_id=?  and create_time BETWEEN ? and ?", "1", merchantId, start, end).Scan(&coinLog)
	} else if storeId == "0" && merchantId == "0" {
		db.SqlDB.Table("coin_log").Select("*").Where("status=? and create_time BETWEEN ? and ?", "1", start, end).Scan(&coinLog)
	}
	for _, val := range coinLog {
		chooseCount = chooseCount + val.ChooseCount
		coinCount = coinCount + val.CoinCount
	}
	//resultMap["gameOrderCount"] = len(coinLog)
	resultMap["chooseCount"] = chooseCount
	resultMap["coinCount"] = coinCount
	return resultMap
}
func GetCoinLogByTime(storeId string, merchantId string, start string, end string) []model.CoinLog {
	var coinLog []model.CoinLog
	if storeId != "0" && merchantId != "0" {
		db.SqlDB.Table("coin_log").Select("*").Where("status=? and merchant_id=? and store_id =? and create_time BETWEEN ? and ?", "1", merchantId, storeId, start, end).Scan(&coinLog)
	} else if storeId == "0" && merchantId != "0" {
		db.SqlDB.Table("coin_log").Select("*").Where("status=? and merchant_id=?  and create_time BETWEEN ? and ?", "1", merchantId, start, end).Scan(&coinLog)
	} else if storeId == "0" && merchantId == "0" {
		db.SqlDB.Table("coin_log").Select("*").Where("status=?  and create_time BETWEEN ? and ?", "1", start, end).Scan(&coinLog)
	}
	return coinLog
}
func GetVisitCount(storeId string, merchantId string, start string, end string) []model.MemberVisitLog {
	var memberVisitLog []model.MemberVisitLog
	loc, _ := time.LoadLocation("Local") //重要：获取时区
	tstart, _ := time.ParseInLocation("2006-01-02 15:04:05", start, loc)
	tend, _ := time.ParseInLocation("2006-01-02 15:04:05", end, loc)
	ts := tstart.Format("2006-01-02")
	te := tend.Format("2006-01-02")
	if storeId != "0" && merchantId != "0" {
		db.SqlDB.Table("member_visit_log").Where("merchant_id = ? and store_id = ? and time_day between ? and ?", merchantId, storeId, ts, te).Group("member_id,time_day").Scan(&memberVisitLog)
	} else if storeId == "0" && merchantId != "0" {
		db.SqlDB.Table("member_visit_log").Where("merchant_id = ? and time_day between ? and ?", merchantId, ts, te).Group("member_id,time_day").Scan(&memberVisitLog)
	} else if storeId == "0" && merchantId == "0" {
		db.SqlDB.Table("member_visit_log").Where("time_day between ? and ?", ts, te).Group("member_id,time_day").Scan(&memberVisitLog)
	}
	return memberVisitLog
}

func GetCoinLogGroupByForMember(storeId string, merchantId string, start string, end string) []model.CoinLog {
	var coinLog []model.CoinLog
	if storeId != "0" && merchantId != "0" {
		db.SqlDB.Table("coin_log").Where("status = ? and store_id = ? and merchant_id = ? and create_time between ? and ?", "1", storeId, merchantId, start, end).Group("member_id").Scan(&coinLog)
	} else if storeId == "0" && merchantId != "0" {
		db.SqlDB.Table("coin_log").Where("status = ?  and merchant_id = ? and create_time between ? and ?", "1", merchantId, start, end).Group("member_id").Scan(&coinLog)
	} else if storeId == "0" && merchantId == "0" {
		db.SqlDB.Table("coin_log").Where("status = ?  and create_time between ? and ?", "1", start, end).Group("member_id").Scan(&coinLog)
	}
	return coinLog
}

func GetCoinStatistice(storeId string, merchantId string, start string, end string, top int) []map[string]interface{} {
	var coinLog []model.CoinLog
	var resultMap []map[string]interface{}
	if storeId != "0" && merchantId != "0" {
		db.SqlDB.Table("coin_log").Select("coin_count,count(*) count").Where("merchant_id = ? and store_id = ? and create_time between ? and ? ", merchantId, storeId, start, end).Group("coin_count").Scan(&coinLog)
	} else if storeId == "0" && merchantId != "0" {
		db.SqlDB.Table("coin_log").Select("coin_count,count(*) count").Where("merchant_id = ?  and create_time between ? and ? ", merchantId, start, end).Group("coin_count").Scan(&coinLog)
	} else if storeId == "0" && merchantId == "0" {
		db.SqlDB.Table("coin_log").Select("coin_count,count(*) count").Where("create_time between ? and ? ", start, end).Group("coin_count").Scan(&coinLog)
	}
	topNum := top
	if topNum > len(coinLog) {
		for _, val := range coinLog {
			result := make(map[string]interface{})
			result["coinCount"] = val.CoinCount
			result["count"] = val.Count
			resultMap = append(resultMap, result)
		}
	} else {
		for i := 0; i < topNum; i++ {
			result := make(map[string]interface{})
			result["coinCount"] = coinLog[i].CoinCount
			result["count"] = coinLog[i].Count
			resultMap = append(resultMap, result)
		}
		var count int64
		count = 0
		for i := topNum; i < len(coinLog); i++ {
			count = count + coinLog[i].Count
		}
		result := make(map[string]interface{})
		result["coinCount"] = "other"
		result["count"] = count
		resultMap = append(resultMap, result)
	}
	return resultMap
}
func getLgzCoinByTime(start string, end string) []model.LogLgzCoin {
	var logLgzCoin []model.LogLgzCoin
	db.SqlDB.Table("log_lgz_coin").Select("*").Where("create_time BETWEEN ? and ?", start, end).Scan(&logLgzCoin)
	return logLgzCoin
}
func getWeiMaQiCoinByTime(start string, end string) []model.LogWeimaqiCoin {
	var logWeimaqiCoin []model.LogWeimaqiCoin
	db.SqlDB.Table("log_weimaqi_coin").Select("*").Where("create_time BETWEEN ? and ?", start, end).Scan(&logWeimaqiCoin)
	return logWeimaqiCoin
}
func getHistoryMember(storeId string, merchantId string, start string, end string) []model.CoinLog {
	var coinLog []model.CoinLog
	if storeId != "0" && merchantId != "0" {
		db.SqlDB.Table("coin_log").Where("status = ? and store_id = ? and merchant_id = ? and create_time between ? and ?", "1", storeId, merchantId, start, end).Scan(&coinLog)
	} else if storeId == "0" && merchantId != "0" {
		db.SqlDB.Table("coin_log").Where("status = ?  and merchant_id = ? and create_time between ? and ?", "1", merchantId, start, end).Scan(&coinLog)
	} else if storeId == "0" && merchantId == "0" {
		db.SqlDB.Table("coin_log").Where("status = ?  and create_time between ? and ?", "1", start, end).Scan(&coinLog)
	}
	return coinLog
}
func GetHistory(storeId string, merchantId string, start string, end string) []map[string]interface{} {
	var resultMap []map[string]interface{}
	loc, _ := time.LoadLocation("Local") //重要：获取时区
	member := getHistoryMember(storeId, merchantId, start, end)
	coinLog := GetCoinLogByTime(storeId, merchantId, start, end)
	//lgzCoin :=getLgzCoinByTime(start,end)
	//weiMaQiCoin :=getWeiMaQiCoinByTime(start,end)
	tstart, err := time.ParseInLocation("2006-01-02 15:04:05", start, loc)
	if err != nil {
		fmt.Println(err)
	}
	tend, _ := time.ParseInLocation("2006-01-02 15:04:05", end, loc)

	ts := tstart.Unix()
	te := tend.Unix()
	var i int64 = 0
	temps := tstart.Format("2006-01-02")
	tempe := tend.Format("2006-01-02")
	var sec int64 = 1800
	if temps == tempe {
		sec = 1800
	} else {
		sec = 60 * 60 * 24
	}

	for i = ts; i < te; i = i + sec {
		var memberCount int = 0
		var onlinegameOrderCount int = 0
		var onlinechooseCount int = 0
		var onlinecoinCount int = 0
		//var offCoinCount int =0
		startTime := time.Unix(i, 0).Format("2006-01-02 15:04:05")
		//endTime := time.Unix(i+1800,0).Format("2006-01-02 15:04:05")
		var memberId map[int64]interface{}
		memberId = make(map[int64]interface{})
		for _, val := range member {
			temp, _ := time.ParseInLocation("2006-01-02 15:04:05", val.CreateTime.String(), loc)
			tempSec := temp.Unix()
			if tempSec >= i && tempSec <= i+sec {
				if _, ok := memberId[val.MemberId]; ok {
					continue
				} else {
					memberId[val.MemberId] = val.MemberId
					memberCount = memberCount + 1
				}
			}
		}
		for _, val := range coinLog {
			temp, _ := time.ParseInLocation("2006-01-02 15:04:05", val.CreateTime.String(), loc)
			tempSec := temp.Unix()
			if tempSec >= i && tempSec <= i+sec {
				onlinegameOrderCount = onlinegameOrderCount + 1
				onlinechooseCount = onlinechooseCount + int(val.ChooseCount)
				onlinecoinCount = onlinecoinCount + int(val.CoinCount)
			}
		}

		//for _,val :=range lgzCoin {
		//	temp, _ := time.ParseInLocation("2006-01-02 15:04:05", val.CreateTime.String(),loc)
		//	tempSec := temp.Unix()
		//	if tempSec >=i && tempSec<=i+sec{
		//		offCoinCount=offCoinCount+int(val.Coin)
		//	}
		//}
		//for _,val := range weiMaQiCoin {
		//	temp, _ := time.ParseInLocation("2006-01-02 15:04:05", val.CreateTime.String(),loc)
		//	tempSec := temp.Unix()
		//	if tempSec >=i && tempSec<=i+sec{
		//		offCoinCount=offCoinCount+int(val.Coinin)
		//	}
		//}

		result := make(map[string]interface{})
		result["time"] = startTime
		result["online_member_count"] = memberCount
		result["online_game_order_count"] = onlinegameOrderCount
		result["online_choose_count"] = onlinechooseCount
		result["online_coin_count"] = onlinecoinCount
		//result["off_coin_count"]=offCoinCount
		resultMap = append(resultMap, result)
	}
	return resultMap
}

func getDeviceFallCountByStoreId(storeId string, merchantId string, start string, end string) int64 {
	var logDeviceFall []model.LogDeviceFall
	//db :=db.SqlDB.Preload("RelDeviceFallRftags")
	if storeId != "0" && merchantId != "0" {
		db.SqlDB.Preload("RelDeviceFallRftags").Find(&logDeviceFall, "merchant_id=? and store_id=? and create_time BETWEEN ? and ?", merchantId, storeId, start, end)
	} else if storeId == "0" && merchantId != "0" {
		db.SqlDB.Preload("RelDeviceFallRftags").Find(&logDeviceFall, "merchant_id=? and create_time BETWEEN ? and ?", merchantId, start, end)
	} else if storeId == "0" && merchantId == "0" {
		db.SqlDB.Preload("RelDeviceFallRftags").Find(&logDeviceFall, "create_time BETWEEN ? and ?", start, end)
	}
	var fallCount int64
	fallCount = 0
	for _, ldfval := range logDeviceFall {
		fallCount = fallCount + 1
		if len(ldfval.RelDeviceFallRftags) > 0 {
			fallCount = fallCount + int64(len(ldfval.RelDeviceFallRftags)) - 1
		}
	}
	return fallCount
}

func GetDeviceList(storeId string, start string, end string, merchantId string) ([]map[string]interface{}, map[string]interface{}) {
	var device []model.Device
	//var   result [] model.Device
	var logWeimaqiCoin []model.LogWeimaqiCoin
	var logLgzCoin []model.LogLgzCoin
	var coinLog []model.CoinLog
	var logDeviceFall []model.LogDeviceFall
	var gameOnLineCount int64 = 0
	var gameOffLineCount int64 = 0
	var chooseCount int64
	var coinCount int64
	var wmqCoin int64
	var lgzCoin int64
	var fallCount int64
	var allOffLineCoin int64
	var allOnlineCoin int64
	var gameOrderCount int64
	online := 0
	allOffLineCoin = 0
	allOnlineCoin = 0
	if storeId != "0" && merchantId != "0" {
		db.SqlDB.Preload("DeviceModular", "modular_type=1").Preload("DeviceModular.CommWeimaqi").Find(&device, "store_id=?", storeId)
		db.SqlDB.Table("log_weimaqi_coin").Select("sum(coinin) coinin,weimaqi_id").Where("create_time BETWEEN ? and ?", start, end).Group("weimaqi_id").Scan(&logWeimaqiCoin)
		db.SqlDB.Table("log_lgz_coin").Select("sum(coin) coin,device_id").Where("create_time BETWEEN ? and ?", start, end).Group("device_id").Scan(&logLgzCoin)
		db.SqlDB.Table("coin_log").Select("*").Where("status=? and merchant_id=? and store_id=? and create_time BETWEEN ? and ?", "1", merchantId, storeId, start, end).Scan(&coinLog)
		db.SqlDB.Preload("RelDeviceFallRftags").Find(&logDeviceFall, "merchant_id=? and store_id=? and create_time BETWEEN ? and ?", merchantId, storeId, start, end)
	} else if storeId == "0" && merchantId != "0" {
		db.SqlDB.Preload("DeviceModular", "modular_type=1").Preload("DeviceModular.CommWeimaqi").Find(&device)
		db.SqlDB.Table("log_weimaqi_coin").Select("sum(coinin) coinin,weimaqi_id").Where("create_time BETWEEN ? and ?", start, end).Group("weimaqi_id").Scan(&logWeimaqiCoin)
		db.SqlDB.Table("log_lgz_coin").Select("sum(coin) coin,device_id").Where("create_time BETWEEN ? and ?", start, end).Group("device_id").Scan(&logLgzCoin)
		db.SqlDB.Table("coin_log").Select("*").Where("status=? and merchant_id=? and create_time BETWEEN ? and ?", "1", merchantId, start, end).Scan(&coinLog)
		db.SqlDB.Preload("RelDeviceFallRftags").Find(&logDeviceFall, "merchant_id=? and create_time BETWEEN ? and ? ", merchantId, start, end)
	} else if storeId == "0" && merchantId == "0" {
		db.SqlDB.Preload("DeviceModular", "modular_type=1").Preload("DeviceModular.CommWeimaqi").Find(&device)
		db.SqlDB.Table("log_weimaqi_coin").Select("sum(coinin) coinin,weimaqi_id").Where("create_time BETWEEN ? and ?", start, end).Group("weimaqi_id").Scan(&logWeimaqiCoin)
		db.SqlDB.Table("log_lgz_coin").Select("sum(coin) coin,device_id").Where("create_time BETWEEN ? and ?", start, end).Group("device_id").Scan(&logLgzCoin)
		db.SqlDB.Table("coin_log").Select("*").Where("status=? and create_time BETWEEN ? and ?", "1", start, end).Scan(&coinLog)
		db.SqlDB.Preload("RelDeviceFallRftags").Find(&logDeviceFall, "create_time BETWEEN ? and ? ", start, end)
	}

	var deviceMaps []map[string]interface{}
	for _, val := range device {
		deviceMap := make(map[string]interface{})
		chooseCount = 0
		coinCount = 0
		wmqCoin = 0
		lgzCoin = 0
		fallCount = 0
		gameOnLineCount = 0
		gameOffLineCount = 0
		if val.Status == 1 {
			online = online + 1 //在线设备
		}
		if len(val.DeviceModular) != 0 {
			//有维码器的话获取每台设备的维码器投币数
			for _, mdval := range val.DeviceModular {
				for _, lwcval := range logWeimaqiCoin {
					if mdval.CommWeimaqi.WeimaqiId == lwcval.WeimaqiId {
						wmqCoin = lwcval.Coinin
						allOffLineCoin = allOffLineCoin + lwcval.Coinin
						gameOffLineCount = gameOffLineCount + 1
					}
				}
			}
		}
		for _, llcval := range logLgzCoin {
			//每台设备的了关注投币数
			if val.ID == llcval.DeviceId {
				lgzCoin = llcval.Coin
				allOffLineCoin = allOffLineCoin + llcval.Coin
				gameOffLineCount = gameOffLineCount + 1
			}
		}

		for _, clval := range coinLog {
			//每台设备线上投币总数等
			if val.ID == clval.DeviceId {
				chooseCount = chooseCount + clval.ChooseCount
				coinCount = coinCount + clval.CoinCount
				allOnlineCoin = allOnlineCoin + clval.CoinCount
				gameOnLineCount = gameOnLineCount + 1
			}
		}

		for _, ldfval := range logDeviceFall {
			if val.StoreId == ldfval.StoreId {
				fallCount = fallCount + 1
				if len(ldfval.RelDeviceFallRftags) > 0 {
					fallCount = fallCount + int64(len(ldfval.RelDeviceFallRftags)) - 1
				}
			}
		}
		deviceMap["online_coin_count"] = coinCount
		deviceMap["offline_coin_count"] = wmqCoin + lgzCoin
		deviceMap["choose_count"] = chooseCount
		//deviceMap["game_order_count"]=len(coinLog)
		deviceMap["device_fall_count"] = fallCount
		deviceMap["id"] = val.ID
		deviceMap["status"] = val.Status
		deviceMap["device_number"] = val.DeviceNumber
		deviceMap["alias"] = val.Alias
		deviceMap["game_online_Count"] = gameOnLineCount
		deviceMap["game_offline_count"] = gameOffLineCount
		gameOrderCount = gameOrderCount + gameOnLineCount
		deviceMaps = append(deviceMaps, deviceMap)
	}
	dataMap := make(map[string]interface{})
	dataMap["online"] = online
	dataMap["allOnlineCoin"] = allOnlineCoin
	dataMap["allOffLineCoin"] = allOffLineCoin
	dataMap["gameOrderCount"] = gameOrderCount
	return deviceMaps, dataMap
}
