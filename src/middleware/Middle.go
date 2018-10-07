package middleware

import (
	"github.com/gin-gonic/gin"
	"SaasServiceGo/src/webgo"
	"strings"
	"github.com/dgrijalva/jwt-go"
	"SaasServiceGo/src/db"
)


//验证
func Middle(c *gin.Context) {
	defer webgo.TryCatchWeb(c)
	tokens := c.GetHeader("authorization")
	NeedLogin := webgo.NEEDLOGIN
	if tokens == "" {
		c.Abort()
		webgo.Result(c, NeedLogin, "no Authorization", nil, nil)
		webgo.Debug("tokens为空：%s","检查token")
		return
	}
	mySignKeyBytes := []byte(webgo.MyKey) //token的Key值
	temp := strings.Split(tokens, " ")
	//验证token
	parseAuth, err := jwt.Parse(temp[1], func(*jwt.Token) (interface{}, error) {
		return mySignKeyBytes, nil
	})
	if err != nil {
		c.Abort()
		webgo.Result(c, NeedLogin, "签名错误", nil, nil)
		webgo.Debug("签名出错：%s",err)
		return
	}
	claim := parseAuth.Claims.(jwt.MapClaims)
	var parmMap map[string]interface{}
	parmMap = make(map[string]interface{})
	for key, val := range claim {
		parmMap[key] = val
	}
	data := parmMap["data"].(map[string]interface{})
	id := webgo.GetResult(data["id"])
	redis := db.CfRedis[2]
	nonce := redis.Get(id).Val()
	//nonce = "lbD0dla1dj8HGwvN3QoRwNqPORfBcqkR"
	if nonce != webgo.GetResult(parmMap["nonceStr"]) {
		webgo.Result(c, NeedLogin, "token err", nil, nil)
		webgo.Debug("nonce不相等：%s","token err")
		c.Abort()
		return
	}
	c.Set("accountId", data["id"])
	c.Set("account", data["account"])
	c.Set("role", data["role"])
	c.Set("merchantId", data["merchant_id"])
	c.Set("storeId", data["store_id"])
	c.Next()
}
