package service

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
	"SaasServiceGo/src/model"
	"SaasServiceGo/src/webgo"
	"SaasServiceGo/src/db"
	"github.com/astaxie/beego"
	"github.com/pkg/errors"
	"fmt"
)

type NewReportFormService struct {

}
type Count struct {
	Count   int64
	Sum     int64
	minTime time.Time
	maxTime time.Time
}

func (ctrl *NewReportFormService)Test(parm string) []model.Device {
	var device []model.Device
	cond, vals, err := whereBuild(map[string]interface{}{
		"status": "1",
		"id in": []int{84, 85, 86,87,88,90},
	})
	if err != nil {
		fmt.Println(err)
	}
	db.SqlDB.Where(cond, vals...).Find(&device)
	return device
}


//根据商户ID和门店ID查找对应数据 判断商户对应门店是否正确
func (ctrl *NewReportFormService) GetRelMerchantStore(parm map[string]interface{}) bool {
	storeId := webgo.GetResult(parm["storeId"])
	merchantId := webgo.GetResult(parm["merchantId"])
	var merchantStore model.MerchantStore
	row := db.SqlDB.Where("store_id=? and merchant_id= ?", storeId, merchantId).First(&merchantStore)
	if row.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

//获取所有门店
func GetAllStore() []model.Store {
	var store []model.Store
	db.SqlDB.Find(&store)
	return store
}

//获取设备报表信息
func (ctrl *NewReportFormService) GetReportDevice(startTime string, endTime string, storeId string, merchantId string,pageModel webgo.PageModel,sort string,status []string,searchParm map[string]string) (map[string]interface{},webgo.PageModel ){
	resultMap := make(map[string]interface{})
	var device []model.Device
	reportDevice := []model.ReportDevice{}
	whereMap :=make(map[string]interface{})
	dWhereMap :=make(map[string]interface{})
	//如果门店不为0查门店 否则看商户是否为0,不为0则查看商户下所有门店的信息，否则查看所有
	rd :=GetDbPageData(pageModel).Table("report_device r").Joins("left join device d on r.device_id=d.id").Select("r.id,sum(device_fall_count)device_fall_count," +
		"sum(choose_count)choose_count,sum(lgz_count)lgz_count," + "sum(weimaqi_count)weimaqi_count,sum(real_coin_count)real_coin_count,sum(coin_count)coin_count,d.alias alias," +
		"d.device_number device_number,d.id device_id,d.store_id store_id,r.create_time create_time,d.status status").Where("r.create_time between ? and ? ",startTime,endTime)
	if status[0]!="" {
		whereMap["d.status in"]=status
		dWhereMap["status in"]=status
	}
	if _, ok := searchParm["deviceNumber"]; ok {
		whereMap["d.device_number like"]="%"+searchParm["deviceNumber"]+"%"
		dWhereMap["device_number like"]="%"+searchParm["deviceNumber"]+"%"
	}
	if _, ok := searchParm["alias"]; ok {
		whereMap["d.alias like"]="%"+searchParm["alias"]+"%"
		dWhereMap["alias like"]="%"+searchParm["alias"]+"%"
	}
	if storeId != "all" {
		whereMap["d.store_id"]=storeId
		dWhereMap["store_id"]=storeId
	} else if merchantId != "all" {
		var merchantStores []model.MerchantStore
		db.SqlDB.Find(&merchantStores, "merchant_id = ?", merchantId)
		var storeIds []int64
		for _, val := range merchantStores {
			storeIds = append(storeIds, val.StoreId)
		}
		whereMap["d.store_id in"]=storeIds
		dWhereMap["store_id in"]=storeIds
	} else {
		var merchantStores []model.MerchantStore
		db.SqlDB.Find(&merchantStores)
		var storeIds []int64
		for _, val := range merchantStores {
			storeIds = append(storeIds, val.StoreId)
		}
		whereMap["d.store_id in"]=storeIds
		dWhereMap["store_id in"]=storeIds
	}
	cond, vals, err := whereBuild(whereMap)
	if err!=nil {
		panic(err)
	}
	rd.Where(cond,vals...).Group("device_id").Order("d."+sort).Scan(&reportDevice)
	dcond, dvals, err := whereBuild(dWhereMap)
	if err!=nil {
		panic(err)
	}
	db.SqlDB.Where(dcond,dvals...).Find(&device)
	var Count Count
	db.SqlDB.Raw("SELECT * FROM (select COUNT(*)count FROM device )s ,(SELECT COUNT(*)sum FROM device WHERE `status`=1)o").Scan(&Count)
	var Total int = len(device)
	pageModel.Total=Total
	resultMap["devices"] = reportDevice
	resultMap["online_device"] = Count.Sum
	resultMap["all_device"] = Count.Count
	return resultMap,pageModel
}

//获取报表信息
func (ctrl *NewReportFormService) GetReportForm(startTime string, endTime string, storeId string, merchantId string,sort string) map[string]interface{} {
	resultMap := make(map[string]interface{})
	dayMemberMap := make(map[string]int64)
	visitMap := make(map[string]int64)
	var reportForm []model.ReportForm
	var newMemberCount Count
	var payMoney Count
	//新增会员数跟充值金额
	db.SqlDB.Table("member").Select("count(*) count").Where("create_time BETWEEN ? and ?", startTime, endTime).Scan(&newMemberCount) //新增会员数
	podb :=db.SqlDB.Table("pay_order").Select("sum(amount) sum,count(*)count").Where("create_time BETWEEN ? and ? ", startTime, endTime)   //充值金额
	if  merchantId!="all"{
		podb=podb.Where("remarks= ? ",merchantId)
	}
	podb.Scan(&payMoney)
	if storeId != "all" {
		db.SqlDB.Table("report_form").Select("*").Where("create_time BETWEEN ? and ? and store_id = ?", startTime, endTime, storeId).Order(sort).Scan(&reportForm)
	} else if merchantId != "all" {
		var merchantStores []model.MerchantStore
		db.SqlDB.Find(&merchantStores, "merchant_id = ?", merchantId)
		var storeIds []int64
		for _, val := range merchantStores {
			storeIds = append(storeIds, val.StoreId)
		}
		db.SqlDB.Table("report_form").Select("*").Where("create_time BETWEEN ? and ? and store_id in (?)", startTime, endTime, storeIds).Order(sort).Scan(&reportForm)
	} else {
		var merchantStores []model.MerchantStore
		db.SqlDB.Table("merchant_store").Select("*").Scan(&merchantStores)
		//	db.SqlDB.Find(&merchantStores)
		var storeIds []int64
		for _, val := range merchantStores {
			storeIds = append(storeIds, val.StoreId)
		}
		db.SqlDB.Table("report_form").Select("*").Where("create_time BETWEEN ? and ? and store_id in (?)", startTime, endTime, storeIds).Order(sort).Scan(&reportForm)
	}
	lineCharts := []map[string]interface{}{}
	var chooseCount int64 = 0
	var deviceFallCount int64 = 0
	var gameOrderCount int64 = 0
	var dayMemberCount int64 = 0
	var offCoinCount int64 = 0
	var onLineCount int64 = 0
	var visitCount int64 = 0
	var realCoinCount int64 = 0
	var mCoinCount int64 = 0
	var coinCount int64 = 0
	resutSpecMap := make(map[int]int)
	otherSpecMap := make(map[string]int)

	//同一天
	if webgo.GetDayIsEqual(startTime, endTime) == true {
		for _, val := range reportForm {
			var lineChar map[string]interface{}
			lineChar = make(map[string]interface{})
			lineChar["time"] = webgo.GetStringDateTime(val.CreateTime)             //时间
			appendFlag :=false
			for i,lineval := range lineCharts {
				if v, ok := lineval["time"]; ok {
					if v==lineChar["time"]{
						lineChar["member_count"] = val.MemberCount+lineval["member_count"].(int64)                              //这段时间的游戏人数
						lineChar["choose_count"] = val.ChooseCount+lineval["choose_count"].(int64)                              //这段时间的游戏次数
						lineChar["online_coin_count"] = val.RealCoinCount+lineval["online_coin_count"].(int64)                        //这段时间的线上投币
						lineChar["offline_coin_count"] = val.WeimaqiCoinCount + val.LgzCoinCount +lineval["offline_coin_count"].(int64)//这段时间的线下
						lineChar["game_order_count"] = val.GameOrderCount +lineval["game_order_count"].(int64)                       //这段时间的订单数量
						lineChar["real_coin_count"] = val.RealCoinCount +lineval["real_coin_count"].(int64)                          //这段时间的真实投币次数
						lineChar["m_coin_count"] = val.MCoinCount    +lineval["m_coin_count"].(int64)                              //这段时间消耗的商户币
						lineChar["coin_count"] = val.CoinCount +lineval["coin_count"].(int64)
						lineCharts[i]=lineChar
						appendFlag=true
					}
				}
			}
			if !appendFlag {
				lineChar["time"] = webgo.GetStringDateTime(val.CreateTime)             //时间
				lineChar["member_count"] = val.MemberCount                              //这段时间的游戏人数
				lineChar["choose_count"] = val.ChooseCount                              //这段时间的游戏次数
				lineChar["online_coin_count"] = val.RealCoinCount                        //这段时间的线上投币
				lineChar["offline_coin_count"] = val.WeimaqiCoinCount + val.LgzCoinCount //这段时间的线下
				lineChar["game_order_count"] = val.GameOrderCount                        //这段时间的订单数量
				lineChar["real_coin_count"] = val.RealCoinCount                          //这段时间的真实投币次数
				lineChar["m_coin_count"] = val.MCoinCount                                //这段时间消耗的商户币
				lineChar["coin_count"] = val.CoinCount
				lineCharts = append(lineCharts, lineChar)
			}
			chooseCount = chooseCount + val.ChooseCount                           //游戏次数
			deviceFallCount = deviceFallCount + val.DeviceFallCount               //出物数
			gameOrderCount = gameOrderCount + val.GameOrderCount                  //订单数
			dayMemberCount = val.DayMemberCount                                   //一整天的游戏人数
			offCoinCount = offCoinCount + val.LgzCoinCount + val.WeimaqiCoinCount //线下投币数
			onLineCount = onLineCount + val.RealCoinCount                         //线上投币数
			visitCount = val.VisitCount                                           //访问人数
			realCoinCount = realCoinCount + val.RealCoinCount                     //真实投币
			mCoinCount = mCoinCount + val.MCoinCount
			coinCount = coinCount + val.CoinCount

			//计算投币规格
			if val.CoinSpec == "" {
				continue
			}
			specMap := webgo.PaserStringToMaps(val.CoinSpec)
			for _, specval := range specMap {
				specOrder, _ := strconv.Atoi(webgo.GetResult(specval["count"]))
				if webgo.GetResult(specval["coinCount"]) != "other" {
					tc, _ := strconv.Atoi(webgo.GetResult(specval["coinCount"]))
					if _, ok := resutSpecMap[tc]; ok {
						resutSpecMap[tc] = resutSpecMap[tc] + specOrder
					} else {
						resutSpecMap[tc] = specOrder
					}
				} else {
					if _, ok := specval["other"]; ok {
						specOther, _ := strconv.Atoi(webgo.GetResult(specval["count"]))
						if _, ok := otherSpecMap["other"]; ok {
							otherSpecMap["other"] = otherSpecMap["other"] + specOther
						} else {
							otherSpecMap["other"] = specOther
						}
					}
				}
			}
		}
		resultMap["choose_count"] = chooseCount
		resultMap["device_fall_count"] = deviceFallCount
		resultMap["game_order_count"] = gameOrderCount
		resultMap["day_member_count"] = dayMemberCount
		resultMap["offcoin_count"] = offCoinCount
		resultMap["online_Count"] = onLineCount
		resultMap["visit_count"] = visitCount
		resultMap["line_chart"] = lineCharts
		resultMap["real_coin_count"] = realCoinCount
		resultMap["m_count_count"] = mCoinCount
		resultMap["coin_count"] = coinCount
		var pct float64 = 0
		if dayMemberCount != 0 {
			pct = float64(coinCount) / float64(dayMemberCount)
		} else {
			pct = 0
		}
		var upt float64 = 0
		if dayMemberCount != 0 {
			upt = float64(gameOrderCount) / float64(dayMemberCount)
		} else {
			upt = 0
		}
		resultMap["pct"] = pct
		resultMap["upt"] = upt
	} else {
		//不同天
		loc, _ := time.LoadLocation("Local")
		start, _ := time.ParseInLocation(webgo.DateTimeFormate, startTime, loc)
		end, _ := time.ParseInLocation(webgo.DateTimeFormate, endTime, loc)
		day := webgo.TimeSub(end, start) //计算日期差
		stemp := start
		tempMap := make(map[string]map[string]interface{})
		//键用来存跨天的日期
		for i := 0; i <= day; i++ {
			//tempMap= append(tempMap,map[string]interface{})
			tempMap[stemp.Format("2006-01-02")] = make(map[string]interface{})
			addOneDay := stemp.Unix() + 86400 //加一天
			st := time.Unix(addOneDay, 0).Format(webgo.DateTimeFormate)
			stemp, _ = time.ParseInLocation(webgo.DateTimeFormate, st, loc)
		}
		dateIdMap := make(map[int64]string)
		//设置日期和计算投币规格
		for _, val := range reportForm {
			//将ID与日期相关联，用Map存放
			dateTime := val.CreateTime.Format("2006-01-02")
			dateIdMap[val.ID] = dateTime
			//计算投币规格
			if val.CoinSpec == "" {
				continue
			}
			specMap := webgo.PaserStringToMaps(val.CoinSpec)
			for _, specval := range specMap {
				specOrder, _ := strconv.Atoi(webgo.GetResult(specval["count"]))
				if webgo.GetResult(specval["coinCount"]) != "other" {
					tc, _ := strconv.Atoi(webgo.GetResult(specval["coinCount"]))
					if _, ok := resutSpecMap[tc]; ok {
						resutSpecMap[tc] = resutSpecMap[tc] + specOrder
					} else {
						resutSpecMap[tc] = specOrder
					}
				} else {
					if _, ok := specval["other"]; ok {
						specOther, _ := strconv.Atoi(webgo.GetResult(specval["count"]))
						if _, ok := otherSpecMap["other"]; ok {
							otherSpecMap["other"] = otherSpecMap["other"] + specOther
						} else {
							otherSpecMap["other"] = specOther
						}
					}
				}
			}
		}

		var maps []map[string]interface{}
		//遍历要查询的日期
		for k, v := range tempMap {
			//	var memberCount int64=0 //这段时间的游戏人数
			var lineChooseCount int64 = 0      //这段时间的游戏次数
			var lineOnlineCoinCount int64 = 0  //这段时间的线上投币
			var lineOffLineCoinCount int64 = 0 //这段时间的线下投币
			var lineGameOrderCount int64 = 0   //这段时间的订单数量
			var lineRealCoinCount int64 = 0    //这段时间的真实投币次数
			var lineMCoinCount int64 = 0       //这段时间消耗的商户币
			var lineCoinCount int64 = 0
			v["time"] = k
			//遍历结果集，通过日期Map
			for _, val := range reportForm {
				if dateIdMap[val.ID] == k {
					chooseCount = chooseCount + val.ChooseCount             //游戏次数
					deviceFallCount = deviceFallCount + val.DeviceFallCount //出物数
					gameOrderCount = gameOrderCount + val.GameOrderCount    //订单数
					dayMemberMap[k] = val.DayMemberCount
					offCoinCount = offCoinCount + val.LgzCoinCount + val.WeimaqiCoinCount //线下投币数
					onLineCount = onLineCount + val.RealCoinCount                         //线上投币数
					visitMap[k] = val.VisitCount
					realCoinCount = realCoinCount + val.RealCoinCount //真实投币
					mCoinCount = mCoinCount + val.MCoinCount
					coinCount = coinCount + val.CoinCount

					lineChooseCount = lineChooseCount + val.ChooseCount
					lineOnlineCoinCount = lineOnlineCoinCount + val.RealCoinCount
					lineOffLineCoinCount = lineOffLineCoinCount + val.WeimaqiCoinCount + val.LgzCoinCount
					lineGameOrderCount = lineGameOrderCount + val.GameOrderCount
					lineRealCoinCount = lineRealCoinCount + val.RealCoinCount
					lineMCoinCount = lineMCoinCount + val.MCoinCount
					lineCoinCount = lineCoinCount + val.CoinCount
					v["member_count"] = val.DayMemberCount
					v["choose_count"] = lineChooseCount
					v["online_coin_count"] = lineOnlineCoinCount
					v["offline_coin_count"] = lineOffLineCoinCount
					v["game_order_count"] = lineGameOrderCount
					v["real_coin_count"] = lineRealCoinCount
					v["m_coin_count"] = lineMCoinCount
					v["coin_count"] = lineCoinCount
				}
			}
			maps = append(maps, v)
		}
		//计算每天游戏人数
		for _, v := range dayMemberMap {
			dayMemberCount = dayMemberCount + v
		}
		//计算访问人数
		for _, v := range visitMap {
			visitCount = visitCount + v
		}
		resultMap["choose_count"] = chooseCount
		resultMap["device_fall_count"] = deviceFallCount
		resultMap["game_order_count"] = gameOrderCount
		resultMap["day_member_count"] = dayMemberCount
		resultMap["off_coin_count"] = offCoinCount
		resultMap["online_count"] = onLineCount
		resultMap["visit_count"] = visitCount
		resultMap["line_chart"] = lineCharts
		resultMap["real_coin_count"] = realCoinCount
		resultMap["m_count_count"] = mCoinCount
		resultMap["coin_count"] = coinCount
		var pct float64 = 0
		if dayMemberCount != 0 {
			pct = float64(coinCount) / float64(dayMemberCount)
		} else {
			pct = 0
		}
		var upt float64 = 0
		if dayMemberCount != 0 {
			upt = float64(gameOrderCount) / float64(dayMemberCount)
		} else {
			upt = 0
		}
		resultMap["pct"] = pct
		resultMap["upt"] = upt
		resultMap["line_chart"] = maps
	}
	resultMap["new_member"] = newMemberCount.Count
	resultMap["pay_money"] = payMoney.Sum
	resultMap["pay_order"]=payMoney.Count
	specMaps := []map[string]int{}
	//投币规格写入返回
	for k, v := range resutSpecMap {
		tempMap := make(map[string]int)
		tempMap["coin_count"] = k
		tempMap["count"] = v
		specMaps = append(specMaps, tempMap)
	}
	resultMap["coin_spec"] = specMaps
	return resultMap
}

//获取历史要插入的报表数据
func (ctrl *NewReportFormService) InsertReportFormHistory(startTime string) []model.ReportForm{
	reportForms := GetData(startTime)
	return reportForms
}
//报表插入
func (ctrl *NewReportFormService) InsertReportForm(startTime string) {
	reportForms := GetData(startTime)
	BatchInsertForm(reportForms, false)
}

//批量更新设备报表
func BatchUpdateDeviceForm(reportDevice []model.ReportDevice)  {
	ids := []int64{}
	formIdMap :=make(map[int64]map[string]interface{})
	for _,val :=range reportDevice {
		formIdMap[val.ID]=webgo.StructToMap(val)
		ids=append(ids, val.ID)
	}
	var report model.ReportDevice
	columnMap := webgo.GetStructTagJson(&report)
	vals := []interface{}{}
	sqlStr := "Update report_device Set "
	for key,val :=range columnMap {
		if key!="DeviceModular" {
			sqlStr += val + " = CASE id"
		}
		for _,rfval :=range reportDevice{
			if key!="DeviceModular" {
				sqlStr+=" when "+ strconv.FormatInt(rfval.ID,10)+" THEN ?"
				vals=append(vals,formIdMap[rfval.ID][key])
			}
		}
		if key!="DeviceModular" {
			sqlStr += " END,"
		}
	}
	sqlStr=beego.Substr(sqlStr,0, len(sqlStr)-1)
	sqlStr += " where id in (?)"
	vals=append(vals, ids)
	err := db.SqlDB.Exec(sqlStr, vals...).Error
	if err!=nil {
		panic(err)
	}
}
//批量更新报表
func BatchUpdateForm(reportForms []model.ReportForm)  {
	ids := []int64{}
	formIdMap :=make(map[int64]map[string]interface{})
	for _,val :=range reportForms {
		formIdMap[val.ID]=webgo.StructToMap(val)
		ids=append(ids, val.ID)
	}
	var report model.ReportForm
	columnMap := webgo.GetStructTagJson(&report)
	vals := []interface{}{}
	sqlStr := "Update report_form Set "
	for key,val :=range columnMap {
		sqlStr+=val+" = CASE id"
		for _,rfval :=range reportForms{
			sqlStr+=" when "+ strconv.FormatInt(rfval.ID,10)+" THEN ?"
			vals=append(vals,formIdMap[rfval.ID][key])
		}
		sqlStr+=" END,"
	}
	sqlStr=beego.Substr(sqlStr,0, len(sqlStr)-1)
	sqlStr += " where id in (?)"
	vals=append(vals, ids)
	err := db.SqlDB.Exec(sqlStr, vals...).Error
	if err!=nil {
		panic(err)
	}
}

//批量插入报表
func BatchInsertForm(reportForms []model.ReportForm, idFlag bool) {
	var sqlStr string
	vals := []interface{}{}
	var inserts []string
	if idFlag == true {
		sqlStr = "INSERT INTO report_form (id,create_time,choose_count, device_fall_count,game_order_count,coin_count,member_count," +
			"day_member_count,visit_count,lgz_coin_count,weimaqi_coin_count,store_id,real_coin_count,m_coin_count,coin_spec) VALUES "
		const rowSQL = "(?,?, ?, ?, ?, ?, ?, ?, ?, ?,?,?,?,?,?)"
		for _, elem := range reportForms {
			inserts = append(inserts, rowSQL)
			vals = append(vals, elem.ID, elem.CreateTime, elem.ChooseCount, elem.DeviceFallCount, elem.GameOrderCount, elem.CoinCount, elem.MemberCount, elem.DayMemberCount, elem.VisitCount, elem.LgzCoinCount, elem.WeimaqiCoinCount, elem.StoreId, elem.RealCoinCount, elem.MCoinCount, elem.CoinSpec)
		}
	} else {
		sqlStr = "INSERT INTO report_form (create_time,choose_count, device_fall_count,game_order_count,coin_count,member_count," +
			"day_member_count,visit_count,lgz_coin_count,weimaqi_coin_count,store_id,real_coin_count,m_coin_count,coin_spec) VALUES "
		const rowSQL = "(?, ?, ?, ?, ?, ?, ?, ?, ?,?,?,?,?,?)"
		for _, elem := range reportForms {
			inserts = append(inserts, rowSQL)
			vals = append(vals, elem.CreateTime, elem.ChooseCount, elem.DeviceFallCount, elem.GameOrderCount, elem.CoinCount, elem.MemberCount, elem.DayMemberCount, elem.VisitCount, elem.LgzCoinCount, elem.WeimaqiCoinCount, elem.StoreId, elem.RealCoinCount, elem.MCoinCount, elem.CoinSpec)
		}
	}
	sqlStr = sqlStr + strings.Join(inserts, ",")
	tx := db.SqlDB.Begin()
	err := tx.Exec(sqlStr, vals...).Error
	if  err != nil {
		tx.Rollback()
		panic(err)
	}else {
		tx.Commit()
	}
}

//计算报表数据
func GetData(startTime string) []model.ReportForm {
	var reportForms []model.ReportForm
	loc, _ := time.LoadLocation("Local") //重要：获取时区
	start, _ := time.ParseInLocation(webgo.DateTimeFormate, startTime, loc)
	st := start.Unix()
	et := st + 1799
	endTime := time.Unix(et, 0).Format(webgo.DateTimeFormate)
	store := GetAllStore()                                        //所有门店
	CoinLogs := GetCoinLogData(startTime, endTime, "")            //获取时间段的coinlog
	logDeviceFall := getDeviceFallCountByTime(startTime, endTime) //时间段的出物数
	dayMember := GetCoinLogData(startTime, endTime, "day")        //一天的游戏人数（后面计算方法：当天去重，隔天未去重）
	memberVisitLog := GetVisitCountByTime(startTime, endTime)     //时间段的游戏
	device := getWeiMaQiRelatDevice()                             //获取所有设备跟维码器的关系
	//logWeiMaQi :=getNewWeiMaQiCoinByCondition(startTime,endTime)//获取时间段内维码器的记录
	//遍历每个门店
	for _, sval := range store {
		var reportForm model.ReportForm
		var chooseCount int64 = 0
		var coinCount int64 = 0
		var fallCount int64 = 0
		var gameOrderCoin int64 = 0
		var memberCount int64 = 0
		var visitCount int64 = 0
		var dayMemberCount int64 = 0
		var weimaqiIds []string
		var deviceIds []int64 //该门店下的设备ID
		var wmqCoin int64 = 0
		var lgzCoin int64 = 0
		var relCoinCount int64 = 0
		var mCoinCount int64 = 0
		reportForm.StoreId = sval.ID
		//时间段的游戏次数，线上投币数，游戏人数，订单数量
		for _, cval := range CoinLogs {
			if sval.ID == cval.StoreId {
				chooseCount = chooseCount + cval.ChooseCount
				coinCount = coinCount + cval.CoinCount
				memberCount = memberCount + 1
				gameOrderCoin = gameOrderCoin + cval.OrderNumber
				relCoinCount = cval.RealCoinCount + relCoinCount
				mCoinCount = mCoinCount + cval.MCoinCount
			}
		}
		//时间段的出物数
		for _, lval := range logDeviceFall {
			if lval.LogDeviceFallMember.ID!=0 {
				if lval.StoreId == sval.ID {
					fallCount = fallCount + 1
					if len(lval.RelDeviceFallRftags) > 0 {
						fallCount = fallCount + int64(len(lval.RelDeviceFallRftags)) - 1
					}
				}
			}
		}
		//时间段的访问人数
		for _, mval := range memberVisitLog {
			if mval.StoreId == sval.ID {
				visitCount = visitCount + 1
			}
		}

		//一天的游戏人数
		for _, dmval := range dayMember {
			if sval.ID == dmval.StoreId {
				dayMemberCount = dayMemberCount + 1
			}
		}

		//获取该门店下维码器的ID
		for _, dval := range device {
			if sval.ID == dval.StoreId {
				for _, dmval := range dval.DeviceModular {
					weimaqiIds = append(weimaqiIds, dmval.CommWeimaqi.WeimaqiId)
				}
				deviceIds = append(deviceIds, dval.ID)
			}
		}

		logWeiMaQi := getNewWeiMaQiCoinByCondition(startTime, endTime, weimaqiIds) //获取该门店下所有维码器的投币记录
		//维码器投币数
		for _, wmval := range logWeiMaQi {
			wmqCoin = wmqCoin + wmval.Coinin
		}
		//乐关注
		logLgz := getNewLgzCoinByCondition(startTime, endTime, deviceIds)
		for _, lgzval := range logLgz {
			lgzCoin = lgzCoin + lgzval.Coin
		}
		coinSpecMap := GetCoinSpec(strconv.FormatInt(sval.ID, 10), startTime, endTime, 4) //投币规格，默认取前4，其余归为其他
		if len(coinSpecMap) > 0 {
			coinSpec, err := json.Marshal(coinSpecMap)
			if err != nil {
				webgo.Debug("报表 coinSpec 解析出错:%+v\n",errors.Wrap(err,"Service GetData"))
			}
			reportForm.CoinSpec = string(coinSpec)
		} else {
			reportForm.CoinSpec = ""
		}
		reportForm.ChooseCount = chooseCount
		reportForm.CreateTime = start
		reportForm.DeviceFallCount = fallCount
		reportForm.GameOrderCount = gameOrderCoin
		reportForm.CoinCount = coinCount
		reportForm.MemberCount = memberCount
		reportForm.DayMemberCount = dayMemberCount
		reportForm.VisitCount = visitCount
		reportForm.LgzCoinCount = lgzCoin
		reportForm.WeimaqiCoinCount = wmqCoin
		reportForm.RealCoinCount = relCoinCount
		reportForm.MCoinCount = mCoinCount
		reportForms = append(reportForms, reportForm)
	}
	return reportForms
}

//获取最新的一条报表设备信息
func (ctrl *NewReportFormService) GetLastReportDevice() model.ReportDevice {
	var reportDevice model.ReportDevice
	db.SqlDB.Table("report_device").Select("*").Where("create_time=(select max(create_time) from report_device)").Limit(1).Scan(&reportDevice)
	return reportDevice
}

//
func (ctrl *NewReportFormService) GetLastReportForm() model.ReportForm {
	var reportForm model.ReportForm
	db.SqlDB.Table("report_form").Select("*").Where("create_time=(select max(create_time) from report_form)").Limit(1).Scan(&reportForm)
	return reportForm
}

func (ctrl *NewReportFormService) InsertReportDeviceHistory(startTime string) []model.ReportDevice{
	reportDevice := GetDeviceData(startTime)
	return reportDevice
}

func (ctrl *NewReportFormService) InsertReportDevice(startTime string){
	reportDevice := GetDeviceData(startTime)
	BatchInsertDevice(reportDevice, false)
}

func BatchInsertDevice(reportDevice []model.ReportDevice, idFlag bool) {
	vals := []interface{}{}
	var inserts []string
	var sqlStr string
	if idFlag == true {
		sqlStr = "INSERT INTO report_device (id,alias,choose_count,device_fall_count,device_number," +
			"device_id,lgz_count,weimaqi_count,real_coin_count,coin_count,status,create_time,store_id) VALUES "
		const rowSQL = "(?,?, ?, ?, ?, ?, ?, ?, ?, ?,?,?,?)"
		for _, elem := range reportDevice {
			inserts = append(inserts, rowSQL)
			vals = append(vals, elem.ID, elem.Alias, elem.ChooseCount, elem.DeviceFallCount, elem.DeviceNumber, elem.DeviceId, elem.LgzCount, elem.WeimaqiCount, elem.RealCoinCount, elem.CoinCount, elem.Status, elem.CreateTime, elem.StoreId)
		}
	} else {
		sqlStr = "INSERT INTO report_device (alias,choose_count,device_fall_count,device_number," +
			"device_id,lgz_count,weimaqi_count,real_coin_count,coin_count,status,create_time,store_id) VALUES "
		const rowSQL = "(?, ?, ?, ?, ?, ?, ?, ?, ?,?,?,?)"
		for _, elem := range reportDevice {
			inserts = append(inserts, rowSQL)
			vals = append(vals, elem.Alias, elem.ChooseCount, elem.DeviceFallCount, elem.DeviceNumber,  elem.DeviceId, elem.LgzCount, elem.WeimaqiCount, elem.RealCoinCount, elem.CoinCount, elem.Status, elem.CreateTime, elem.StoreId)
		}
	}
	sqlStr = sqlStr + strings.Join(inserts, ",")
	tx := db.SqlDB.Begin()
	err := tx.Exec(sqlStr, vals...).Error
	if  err != nil {
		tx.Rollback()
		panic(err)
	}else {
		tx.Commit()
	}
}

func GetDeviceData(startTime string) []model.ReportDevice {
	var reportDevices []model.ReportDevice
	loc, _ := time.LoadLocation("Local") //重要：获取时区
	start, _ := time.ParseInLocation(webgo.DateTimeFormate, startTime, loc)
	st := start.Unix()
	et := st + 1799
	endTime := time.Unix(et, 0).Format(webgo.DateTimeFormate)

	devices := getWeiMaQiRelatDevice()
	storeWeiMaQiMap := make(map[int64][]string)
	storeDeviceIdMap := make(map[int64][]int64)
	var logDeviceFall []model.LogDeviceFall
	//出物数
	db.SqlDB.Preload("RelDeviceFallRftags").Find(&logDeviceFall, "create_time BETWEEN ? and ?", startTime, endTime)
	deviceModularMap := make(map[int64][]model.DeviceModular)
	//遍历所有设备
	for _, deviceVal := range devices {
		var reportDevice model.ReportDevice
		reportDevice.StoreId = deviceVal.StoreId
		reportDevice.DeviceId = deviceVal.ID
		deviceModularMap[deviceVal.ID]=deviceVal.DeviceModular
		reportDevice.Alias = deviceVal.Alias
		reportDevice.CreateTime = start
		reportDevice.Status = deviceVal.Status
		reportDevices = append(reportDevices, reportDevice)
		//获取每个门店的所有设备的维码器ID
		if _, ok := storeWeiMaQiMap[deviceVal.StoreId]; ok {
			storeWeiMaQiIds := storeWeiMaQiMap[deviceVal.StoreId]
			for _, dmVal := range deviceVal.DeviceModular {
				storeWeiMaQiIds = append(storeWeiMaQiIds, dmVal.CommWeimaqi.WeimaqiId)
			}
			storeWeiMaQiMap[deviceVal.StoreId] = storeWeiMaQiIds
		} else {
			var storeWeiMaQiIds []string
			for _, dmVal := range deviceVal.DeviceModular {
				storeWeiMaQiIds = append(storeWeiMaQiIds, dmVal.CommWeimaqi.WeimaqiId)
			}
			storeWeiMaQiMap[deviceVal.StoreId] = storeWeiMaQiIds
		}

		//获取每个门店的所有设备的设备ID
		if _, ok := storeDeviceIdMap[deviceVal.StoreId]; ok {
			storeDeviceIds := storeDeviceIdMap[deviceVal.StoreId]
			storeDeviceIds = append(storeDeviceIds, deviceVal.ID)
			storeDeviceIdMap[deviceVal.StoreId] = storeDeviceIds
		} else {
			var storeDeviceIds []int64
			storeDeviceIds = append(storeDeviceIds, deviceVal.ID)
			storeDeviceIdMap[deviceVal.StoreId] = storeDeviceIds
		}
	}
	//按照每个门店去查每个门店所有维码器的投币记录，计算对应设备的维码器投币数
	for _, v := range storeWeiMaQiMap {
		var logWeiMaQi []model.LogWeimaqiCoin
		db.SqlDB.Table("log_weimaqi_coin").Select("sum(coinin) coinin,weimaqi_id").Where("create_time BETWEEN ? and ? and weimaqi_id in (?)", startTime, endTime, v).Group("weimaqi_id").Scan(&logWeiMaQi)
		for i, rval := range reportDevices {
				if  _, ok := deviceModularMap[rval.DeviceId]; ok  {
					for  _, dmval := range deviceModularMap[rval.DeviceId] {
						for _, wmq := range logWeiMaQi {
							if dmval.CommWeimaqi.WeimaqiId == wmq.WeimaqiId {
								reportDevices[i].WeimaqiCount = rval.WeimaqiCount + wmq.Coinin
							}
						}
					}
				}
		}
	}

	//乐关注 跟线上投币
	for _, v := range storeDeviceIdMap {
		var logLgzCoin []model.LogLgzCoin
		var coinLog []model.CoinLog
		db.SqlDB.Table("log_lgz_coin").Select("sum(coin) coin,device_id").Where("create_time BETWEEN ? and ? and device_id in(?)", startTime, endTime, v).Group("device_id").Scan(&logLgzCoin)
		db.SqlDB.Table("coin_log").Select("*").Where("create_time BETWEEN ? and ? and status=? and device_id in (?)", startTime, endTime, "1", v).Scan(&coinLog)

		for i, rval := range reportDevices {
			//乐关注
			for _, lgzVal := range logLgzCoin {
				if rval.DeviceId == lgzVal.DeviceId {
					reportDevices[i].LgzCount = rval.LgzCount + lgzVal.Coin
					//reportDevices[i].GameOfflineCount = reportDevices[i].GameOfflineCount + 1
				}
			}
			//线上投币
			for _, clgVal := range coinLog {
				if clgVal.DeviceId == rval.DeviceId {
					reportDevices[i].CoinCount = reportDevices[i].CoinCount + clgVal.CoinCount
					reportDevices[i].ChooseCount = reportDevices[i].ChooseCount + clgVal.ChooseCount
					reportDevices[i].RealCoinCount = reportDevices[i].RealCoinCount + clgVal.RealCoinCount
					//reportDevices[i].GameOnlineCount = reportDevices[i].GameOnlineCount + 1
				}
			}
		}
	}

	//出物数
	for i, rval := range reportDevices {
		for _, fval := range logDeviceFall {
			if fval.DeviceId == rval.DeviceId {
				reportDevices[i].DeviceFallCount = reportDevices[i].DeviceFallCount + 1
				if len(fval.RelDeviceFallRftags) > 0 {
					reportDevices[i].DeviceFallCount = reportDevices[i].DeviceFallCount + int64(len(fval.RelDeviceFallRftags)) - 1
				}
			}
		}
	}
	return reportDevices
}

//更新报表信息
func (ctrl *NewReportFormService) UpdateRepormForm(startTime string) {
	parmTime := webgo.TimeZone(startTime)
	var oldReportForms []model.ReportForm
	//获取旧的数据
	db.SqlDB.Table("report_form").Select("*").Where("create_time = ? ", parmTime).Scan(&oldReportForms)
	var ids []int64
	if len(oldReportForms) > 0 {
		newReportForm := GetData(parmTime) //获取新想数据
		var tempReportForm []model.ReportForm
		var j int =0
		for i, newVal := range newReportForm {
			for _, oldVal := range oldReportForms {
				if newVal.StoreId == oldVal.StoreId {
					//保留旧数据的ID
					newVal.ID = oldVal.ID
					newReportForm[i].ID = oldVal.ID
					ids = append(ids, oldVal.ID)
					if !(oldVal==newVal) {
						tempReportForm=append(tempReportForm, newVal)
						j=j+1
						if j==10 {
							BatchUpdateForm(tempReportForm)
							tempReportForm=append(tempReportForm[:0],tempReportForm[len(tempReportForm):]...)
							j=0
						}
					}
				}
			}
		}
		if len(tempReportForm)>0 {
			BatchUpdateForm(newReportForm)
		}
	}
}

//更新设备信息
func (ctrl *NewReportFormService) UpdateReportDevice(startTime string) {
	parmTime := webgo.TimeZone(startTime)
	var oldReportDevice []model.ReportDevice
	db.SqlDB.Table("report_device").Select("*").Where("create_time = ? ", parmTime).Scan(&oldReportDevice)
	var ids []int64
	if len(oldReportDevice) > 0 {
		newReportDevice := GetDeviceData(parmTime)
		var tempReportDevice []model.ReportDevice
		var j int=0
		for i, newVal := range newReportDevice {
			for _, oldVal := range oldReportDevice {
				if newVal.DeviceId == oldVal.DeviceId {
					newVal.ID = oldVal.ID
					newReportDevice[i].ID = oldVal.ID
					ids = append(ids, oldVal.ID)
					//新旧数据不相等 更新
					if !(oldVal==newVal){
						tempReportDevice=append(tempReportDevice, newVal)
						j=j+1
						if j==10 {
							BatchUpdateDeviceForm(tempReportDevice)
							tempReportDevice=append(tempReportDevice[:0],tempReportDevice[len(tempReportDevice):]...)
							j=0
						}
					}
				}
			}
		}
		if len(tempReportDevice)>0 {
			BatchUpdateDeviceForm(tempReportDevice)
		}
	}
}

//获取线上投币的信息 flag 为day查询一天的数据
func GetCoinLogData(start string, end string, flag string) []model.CoinLog {
	var coinLog []model.CoinLog
	parmStart := start
	parmEnd := end
	if flag == "day" {
		st := strings.Split(start, " ")
		parmStart = st[0] + " 00:00:00"
		parmEnd = st[0] + " 23:59:59"
	}
	db.SqlDB.Table("coin_log").Select("SUM(choose_count)choose_count,SUM(coin_count)coin_count,COUNT(*)order_number,sum(real_coin_count)real_coin_count,sum(m_coin_count)m_coin_count,store_id,member_id").Where(""+
		"status=? and create_time BETWEEN ? and ?", "1", parmStart, parmEnd).Group("store_id,member_id").Scan(&coinLog)
	return coinLog
}

//获取线上投币的最早的一条记录
func (ctrl *NewReportFormService) GetCoinLogEarly() time.Time {
	var count model.CoinLog
	db.SqlDB.Table("coin_log").Select("MIN(CREATE_time)create_time").Scan(&count)
	t :=count.CreateTime.CTime()
	return t
}

//获取投币记录的最早的一条记录
func (ctrl *NewReportFormService) GetDeviceEarly() time.Time {
	var coinlogCount model.CoinLog
	var lgzCount model.LogLgzCoin
	var wmqCount model.LogWeimaqiCoin
	db.SqlDB.Table("coin_log").Select("MIN(CREATE_time)create_time").Scan(&coinlogCount)
	db.SqlDB.Table("log_lgz_coin").Select("MIN(CREATE_time)create_time").Scan(&lgzCount)
	db.SqlDB.Table("log_weimaqi_coin").Select("MIN(CREATE_time)create_time").Scan(&wmqCount)

	coinlogTime := coinlogCount.CreateTime.CTime()
	lgzTime := lgzCount.CreateTime.CTime()
	weimaqiTime := wmqCount.CreateTime.CTime()
	cl := coinlogTime.Before(lgzTime)
	cw := coinlogTime.Before(weimaqiTime)
	lw := lgzTime.Before(weimaqiTime)
	if cl && cw {
		return coinlogTime
	} else if lw {
		return lgzTime
	} else {
		return weimaqiTime
	}

}

//获取参观人数
func GetVisitCountByTime(start string, end string) []model.MemberVisitLog {
	var memberVisitLog []model.MemberVisitLog
	loc, _ := time.LoadLocation("Local") //重要：获取时区
	tstart, _ := time.ParseInLocation(webgo.DateTimeFormate, start, loc)
	tend, _ := time.ParseInLocation(webgo.DateTimeFormate, end, loc)
	ts := tstart.Format("2006-01-02")
	te := tend.Format("2006-01-02")
	//每半小时查一次 如果不等说明开始时间为23.30 结束时间跨天
	if ts != te {
		te = ts
	}
	db.SqlDB.Table("member_visit_log").Where("time_day between ? and ?", ts, te).Group("member_id,time_day").Scan(&memberVisitLog)
	return memberVisitLog
}

//获取线上投币规格与其对应的订单数
func GetCoinSpec(storeId string, start string, end string, top int) []map[string]interface{} {
	var coinLog []model.CoinLog
	var resultMap []map[string]interface{}
	db.SqlDB.Table("coin_log").Select("coin_count,count(*) count").Where("create_time between ? and ? and store_id= ?", start, end, storeId).Group("coin_count").Scan(&coinLog)
	for _, val := range coinLog {
		result := make(map[string]interface{})
		result["coinCount"] = val.CoinCount
		result["count"] = val.Count
		resultMap = append(resultMap, result)
	}
	return resultMap
}

//获取乐关注日志信息
func getNewLgzCoinByCondition(start string, end string, deviceIds []int64) []model.LogLgzCoin {
	var logLgzCoin []model.LogLgzCoin
	db.SqlDB.Table("log_lgz_coin").Select("sum(coin) coin,device_id").Where("create_time BETWEEN ? and ? and device_id in (?)", start, end, deviceIds).Group("device_id").Scan(&logLgzCoin)
	return logLgzCoin
}

//获取维码器日志信息
func getNewWeiMaQiCoinByCondition(start string, end string, weimaqiIds []string) []model.LogWeimaqiCoin {
	var logWeimaqiCoin []model.LogWeimaqiCoin
	db.SqlDB.Table("log_weimaqi_coin").Select("sum(coinin) coinin,weimaqi_id").Where("create_time BETWEEN ? and ? and weimaqi_id in (?)", start, end, weimaqiIds).Group("weimaqi_id").Scan(&logWeimaqiCoin)
	return logWeimaqiCoin
}

//获取设备箱关联的维码器信息
func getWeiMaQiRelatDevice() []model.Device {
	var device []model.Device
	db.SqlDB.Preload("DeviceModular", "modular_type=1").Preload("DeviceModular.CommWeimaqi").Find(&device)
	return device
}

//获取出物数据
func getDeviceFallCountByTime(start string, end string) []model.LogDeviceFall {
	var logDeviceFall []model.LogDeviceFall
	//db.SqlDB.Preload("RelDeviceFallRftags").Find(&logDeviceFall,"create_time BETWEEN ? and ?")
	db.SqlDB.Preload("RelDeviceFallRftags").Preload("LogDeviceFallMember").Find(&logDeviceFall, "create_time BETWEEN ? and ?", start, end)
	return logDeviceFall
}
