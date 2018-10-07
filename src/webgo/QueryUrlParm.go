package webgo

import (
	"github.com/gin-gonic/gin"
	"strings"
)
//表单查询是否排序
func IsSort(ctx *gin.Context) (string,bool) {
	 parm:=ctx.Query("sort")
	if parm !=""{
		arrParm := strings.Split(parm,",")
		var result string
		for _,val :=range arrParm {
			isDesc :=strings.Contains(val,"-")
			if isDesc ==true {
				val=strings.Replace(val,"-","",-1)
				result=result+val+" desc,"
			}else {
				result=result+val+" asc,"
			}
		}
		result=result[0:len(result)-1]
		return result,true
	}
	return "",false
}
