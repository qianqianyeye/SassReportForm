package webgo

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Result(ctx *gin.Context, status int, msg interface{}, data interface{}, meta interface{}) {
	HttpStatus(ctx, status)
	ctx.JSON(ctx.Writer.Status(), gin.H{"status": status, "msg": msg, "data": data, "meta": meta})
}

func ResultOk(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": data, "msg": ""})
}
func ResultList(ctx *gin.Context, data interface{}, total int64) {
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "rows": data, "msg": "", "total": total})
}
func ResultOkMsg(ctx *gin.Context, data interface{}, msg string) {
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": data, "msg": msg})
}

func ResultFail(ctx *gin.Context, err interface{}) {
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "data": nil, "msg": err})
}

func ResultFailData(ctx *gin.Context, data interface{}, err interface{}) {
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "data": data, "msg": err})
}

func HttpStatus(c *gin.Context, status int) {
	method := c.Request.Method
	if c.Request.URL.String() == "/api/login" {
		c.Status(200)
		return
	}

	methodArr := []string{"PUT", "POST", "PATCH"}
	for _, val := range methodArr {
		if val == method {
			if status == SUCCESS {
				c.Status(201)
				return
			}
		}
	}
	if status == NEEDLOGIN {
		c.Status(401)
		return
	}
	if status == ERROR {
		c.Status(400)
		return
	}
	if status == NOAUTH {
		c.Status(403)
	}
}
