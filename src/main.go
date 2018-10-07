package main

import (
	"github.com/gin-gonic/gin"
	"SaasServiceGo/src/db"
	"time"
	"SaasServiceGo/src/webgo"
	"github.com/itsjamie/gin-cors"
	"SaasServiceGo/src/controller"
	"SaasServiceGo/src/schedule"
	"flag"
)

func registerRouter(router *gin.Engine) {
	new(controller.ReportFormController).Router(router)
//	new(controller.TestController).Router(router)
}
func main() {
	defer webgo.TryCatch()
	//l4g.LoadConfiguration("config/log4g.xml") //使用加载配置文件,类似与java的log4j.propertites
	//defer l4g.Close()               //注:如果不是一直运行的程序,请加上这句话,否则主线程结束后,也不会输出和log到日志文件
	dataBase := flag.Bool("MySql",false,"true :线上，false: 线下 默认:false")
	flag.Parse()
	//*dataBase=true
	db.InitDB(*dataBase) //初始化数据库
	db.InitRedis(*dataBase) //初始化Redis
	defer db.SqlDB.Close()
	schedule.InsertHistory() //计算报表历史数据

	go schedule.Report()     //半小时插入数据
	go schedule.KeepMysql()  //保持与数据库的连接
	router := gin.Default()
	//网页跨域问题
	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET,PUT,POST,DELETE,OPTIONS",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))
	registerRouter(router)
	router.Run(":8030")
}