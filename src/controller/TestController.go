package controller

import (
	"SaasServiceGo/src/webgo"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v8"
	"time"
	"reflect"
	"net/http"
)

type TestController struct {
	webgo.Controller
}

type Booking struct {
	CheckIn  time.Time `form:"check_in" binding:"required,TimeValid" time_format:"2006-01-02"`
	CheckOut time.Time `form:"check_out" binding:"required" time_format:"2006-01-02"`
}

func bookableDate(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	if date, ok := field.Interface().(time.Time); ok {
		today := time.Now()
		if today.Year() > date.Year() || today.YearDay() > date.YearDay() {
			return false
		}
	}
	return true
}

func getBookable(c *gin.Context) {
	var b Booking
	if err := c.ShouldBindWith(&b, binding.Query); err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Booking dates are valid!"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (ctrl *TestController) Router(router *gin.Engine) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("TimeValid", webgo.TimeValid)
	}
	router.POST("/bookable", getBookable)
	//r := router.Group("go/api/test")
	//if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	//	v.RegisterValidation("midtest", midtest)
	//}
	////r.POST("report_form/store/:store_id/coin",ctrl.getCoinLog)
	//r.GET("reportForm/store/:store_id", ctrl.getReportForm)
	//r.POST("reportDevice/store/:store_id", ctrl.getReportDevice)
}


func midtest(v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
)bool {
	if date, ok := field.Interface().(time.Time); ok {
		today := time.Now()
		if today.Year() > date.Year() || today.YearDay() > date.YearDay() {
			return false
		}
	}
	return true
}
