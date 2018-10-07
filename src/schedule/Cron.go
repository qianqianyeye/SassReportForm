package schedule

import (

	"github.com/robfig/cron"
	"time"

	"SaasServiceGo/src/webgo"
	"SaasServiceGo/src/model"
	"SaasServiceGo/src/service"
	"SaasServiceGo/src/db"
)

var newReportFormService service.NewReportFormService

//var startTime string = "2018-07-11 00:00:00"

func Report() {
	c := cron.New()
	//在每小时的0分和30分执行插入数据
	c.AddFunc("0 0,30 * * * *", func() {
		defer webgo.TryCatch()
		start := time.Now()
		//sec := start.Second()
		startUnix := start.Unix()
		//if sec > 0 {
		//	startUnix = startUnix - int64(sec)
		//}
		startTime := time.Unix(startUnix, 0).Format(webgo.DateTimeFormate)
		startTime = webgo.TimeZone(startTime)
		newReportFormService.InsertReportForm(startTime)
		newReportFormService.InsertReportDevice(startTime)
	})
	//每隔5分钟分钟更新一次
	c.AddFunc("0 */5 * * * *", func() {
		defer webgo.TryCatch()
		start := time.Now().Format(webgo.DateTimeFormate)
		//	start :="2018-08-05 01:00:00"
		newReportFormService.UpdateRepormForm(start)
		newReportFormService.UpdateReportDevice(start)
	})
	//在每小时的29分59秒更新0分到29分59秒的数据  在59分59秒更新30分到59分59秒的数据
	c.AddFunc("59 29,59 * * * *", func() {
		defer webgo.TryCatch()
		//fmt.Println("29 59 分更新！")
		st := time.Now().Unix() - 200
		start := time.Unix(st, 0).Format(webgo.DateTimeFormate)
		newReportFormService.UpdateRepormForm(start)
		newReportFormService.UpdateReportDevice(start)
	})
	//c.AddFunc("*/10 * * * * *", func() {
	//	newReportFormService.InsertReportDevice(startTime)
	//	newReportFormService.InsertReportForm(startTime)
	//	loc, _ := time.LoadLocation("Local")                            //重要：获取时区
	//	start, _ := time.ParseInLocation("2006-01-02 15:04:05",startTime,loc)
	//	st :=start.Unix()
	//	var et int64=0
	//	et=st+1800
	//	startTime = time.Unix(et, 0).Format("2006-01-02 15:04:05")
	//})
	c.Start()
}

//将历史数据计算到报表
func InsertHistory() {
	defer webgo.TryCatch()
	webgo.Debug("检查历史数据%s","...")
	start := time.Now().Format(webgo.DateTimeFormate)
	startTime := webgo.TimeZone(start)
	loc, _ := time.LoadLocation("Local") //重要：获取时区
	startT, _ := time.ParseInLocation(webgo.DateTimeFormate, startTime, loc)
	startUnix := startT.Unix()
	//defaultTime := "2018-07-11 00:00:00"
	reportForm := newReportFormService.GetLastReportForm()     //获取最新的时间一条数据
	reportDevice := newReportFormService.GetLastReportDevice() //获取最新的时间一条数据
	var resultForm [][]model.ReportForm
	var resultDevice [][]model.ReportDevice
	formFlag := false
	deviceFlag := false

	if reportForm.ID != 0 {
		flagUnix := reportForm.CreateTime.Unix()
		if flagUnix != startUnix {
			//获取历史时间段内每半小时的数据
			for i := flagUnix + 1800; i <= startUnix; i = i + 1800 {
				Time := time.Unix(i, 0).Format(webgo.DateTimeFormate)
				reportForm :=newReportFormService.InsertReportFormHistory(Time)
				resultForm=append(resultForm, reportForm)
			}
			//批量插入
			var insertForm []model.ReportForm
			var j int =0
			for _,val := range resultForm {
				for _,rval := range val   {
					insertForm=append(insertForm, rval)
					j=j+1
					if j==240 {
						service.BatchInsertForm(insertForm,false)
						insertForm=append(insertForm[:0],insertForm[len(insertForm):]...)
						j=0
					}
				}
			}
			if len(insertForm)>0 {
				service.BatchInsertForm(insertForm,false)
			}
		} else {
			formFlag = true
		}

	} else {
		//loc, _ := time.LoadLocation("Local") //重要：获取时区
		//start, _ := time.ParseInLocation("2006-01-02 15:04:05", defaultTime, loc)
		starts := newReportFormService.GetCoinLogEarly()
		st := starts.Unix()
		for i := st ; i <= startUnix; i = i + 1800 {
			Time := time.Unix(i, 0).Format(webgo.DateTimeFormate)
			reportForm :=newReportFormService.InsertReportFormHistory(Time)
			resultForm=append(resultForm, reportForm)
		}
		var insertForm []model.ReportForm
		var j int =0
		for _,val := range resultForm {
			for _,rval := range val   {
				insertForm=append(insertForm, rval)
				j=j+1
				if j==240 {
					service.BatchInsertForm(insertForm,false)
					insertForm=append(insertForm[:0],insertForm[len(insertForm):]...)
					j=0
				}
			}
		}
		service.BatchInsertForm(insertForm,false)
	}

	if reportDevice.ID != 0 {
		flagUnix := reportDevice.CreateTime.Unix()
		if flagUnix != startUnix {
			for i := flagUnix + 1800; i <= startUnix; i = i + 1800 {
				Time := time.Unix(i, 0).Format(webgo.DateTimeFormate)
				reportDevice:=newReportFormService.InsertReportDeviceHistory(Time)
				resultDevice=append(resultDevice, reportDevice)
			}
			var insertDevice []model.ReportDevice
			var j int =0
			for _,val := range resultDevice {
				for _,rval := range val   {
					insertDevice=append(insertDevice, rval)
					j=j+1
					if j==240{
						service.BatchInsertDevice(insertDevice,false)
						insertDevice=append(insertDevice[:0],insertDevice[len(insertDevice):]...)
						j=0
					}
				}
			}
			if len(insertDevice)>0 {
				service.BatchInsertDevice(insertDevice,false)
			}
		} else {
			deviceFlag = true
		}

	} else {
		starts := newReportFormService.GetDeviceEarly()
		st := starts.Unix()
		//获取历史时间段内每半小时的数据
		for i := st ; i <= startUnix; i = i + 1800 {
			Time := time.Unix(i, 0).Format(webgo.DateTimeFormate)
			reportDevice:=newReportFormService.InsertReportDeviceHistory(Time)
			resultDevice=append(resultDevice, reportDevice)
		}
		var j int =0
		//批量插入
		var insertDevice []model.ReportDevice
		for _,val := range resultDevice {
			for _,rval := range val   {
				insertDevice=append(insertDevice, rval)
				j=j+1
				if j==240{
					service.BatchInsertDevice(insertDevice,false)
					insertDevice=append(insertDevice[:0],insertDevice[len(insertDevice):]...)
					j=0
				}
			}
		}
		if len(insertDevice)>0 {
			service.BatchInsertDevice(insertDevice,false)
		}
	}
	//插入完历史数据在次调用校验，直到数据为最新的
	if formFlag == false || deviceFlag == false {
		InsertHistory()
	}
	webgo.Debug("插入%s!","success")
}

func KeepMysql() {
	c := cron.New()
	c.AddFunc("*/5 * * * * ", func() {
		defer webgo.TryCatch()
		db.SqlDB.DB().Ping()
	})
	c.Start()
}
