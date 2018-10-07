package controller

import (

	"github.com/gin-gonic/gin"
	"SaasServiceGo/src/service"
	"SaasServiceGo/src/webgo"
	"strings"

	"time"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"gopkg.in/go-playground/validator.v8"

	"SaasServiceGo/src/middleware"
)

type ReportFormController struct {
	webgo.Controller
}

type parmValid struct {
	Current int `form:"current"  binding:"required,PageValid"`
	Page_size int `form:"page_size"  binding:"required,PageValid"`
	Start_time time.Time `form:"start_time" binding:"required,TimeValid" time_format:"2006-01-02 15:04:05"`
	End_time time.Time `form:"end_time" binding:"required,TimeValid" time_format:"2006-01-02 15:04:05"`
	Merchant_id string `form:"merchant_id"  binding:"required,MerchantIdAndStoreIdValid"`
	Store_id string `form:"store_id"  binding:"required,MerchantIdAndStoreIdValid"`
}

var newReportService service.NewReportFormService

func (ctrl *ReportFormController) Router(router *gin.Engine) {
	r := router.Group("go/api/v1",middleware.Middle)
	//r := router.Group("go/api/v1")
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("PageValid", webgo.PageValid)
		v.RegisterValidation("TimeValid", webgo.TimeValid)
		v.RegisterValidation("MerchantIdAndStoreIdValid", webgo.MerchantIdAndStoreIdValid)
	}
	r.GET("reportForm/store/:store_id", ctrl.getReportForm)
	r.GET("reportDevice/store/:store_id", ctrl.getReportDevice)
	r.GET("reportDevice/test", ctrl.test)
}
func (ctrl *ReportFormController) test(ctx *gin.Context){
	parm :=ctx.Query("parm")
	data:=newReportService.Test(parm)
	webgo.Result(ctx,0,"",data,"")
}

func (ctrl *ReportFormController) getReportForm(ctx *gin.Context) {
	var parmValid parmValid
	if err := ctx.ShouldBindWith(&parmValid, binding.Query); err == nil {
		storeId := ctx.Query("store_id")
		if ctx.Keys["role"] == "store" {
			if storeId != webgo.GetResult(ctx.Keys["storeId"]) {
				webgo.Result(ctx, webgo.ILLEGALARGUMENT, "无权查看！", nil, nil)
				return
			}
		}
		merchantId := ctx.Query("merchant_id")
		var parm map[string]interface{}
		parm = make(map[string]interface{})
		parm["merchantId"] = merchantId
		parm["storeId"] = storeId
		flag := false
		if storeId != "all" && merchantId !="all"{
			flag = newReportService.GetRelMerchantStore(parm)
		} else {
			flag = true
		}
		if flag {
			sort,isSort :=webgo.IsSort(ctx)
			if isSort {
				result := newReportService.GetReportForm(ctx.Query("start_time"), ctx.Query("end_time"), storeId, webgo.GetResult(merchantId),sort)
				webgo.Result(ctx, webgo.SUCCESS, nil, result, nil)
			}else {
				sort ="create_time"
				result := newReportService.GetReportForm(ctx.Query("start_time"), ctx.Query("end_time"), storeId, webgo.GetResult(merchantId),sort)
				webgo.Result(ctx, webgo.SUCCESS, nil, result, nil)
			}
		} else {
			webgo.Result(ctx, webgo.ILLEGALARGUMENT, "无权查看！", nil, nil)
		}
	} else {
		ctx.JSON(http.StatusOK, gin.H{"error": err.Error()})
	}
}
func (ctrl *ReportFormController) getReportDevice(ctx *gin.Context) {
	storeId := ctx.Query("store_id")
	//storeId := webgo.GetResult(ctx.Param("store_id"))
	pageSize:=ctx.Query("page_size")
	current :=ctx.Query("current")
	pageModel := webgo.GetPageInfo(pageSize,current)
	statusParm :=ctx.Query("status")
	deviceNumber :=ctx.Query("device_number")
	alias :=ctx.Query("alias")
	searchParm := make(map[string]string)

	if deviceNumber!="" {
		searchParm["deviceNumber"]=deviceNumber
	}
	if alias!="" {
		searchParm["alias"]=alias
	}
	status :=strings.Split(statusParm,",")

	if ctx.Keys["role"] == "store" {
		if storeId != webgo.GetResult(ctx.Keys["storeId"]) {
			webgo.Result(ctx, webgo.ILLEGALARGUMENT, "无权查看！", nil, nil)
			return
		}
	}
	merchantId := ctx.Query("merchant_id")
	//var parm map[string]interface{}
	parm := make(map[string]interface{})
	parm["merchantId"] = merchantId
	parm["storeId"] = storeId
	flag := false
	if storeId != "all" && merchantId!="all"{
		flag = newReportService.GetRelMerchantStore(parm)
	} else {
		flag = true
	}
	if flag {
		sort,isSort :=webgo.IsSort(ctx)
		if !isSort {
			sort="create_time"
		}
		result,pageModel:= newReportService.GetReportDevice(ctx.Query("start_time"), ctx.Query("end_time"), storeId, webgo.GetResult(merchantId),pageModel,sort,status,searchParm)
		webgo.Result(ctx, webgo.SUCCESS, nil, result, pageModel)
	} else {
		webgo.Result(ctx, webgo.ILLEGALARGUMENT, "无权查看！", nil, nil)
	}
}
