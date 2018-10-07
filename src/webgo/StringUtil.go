package webgo

import (
	"encoding/json"
	"fmt"
	"strconv"
)

//转换为字符串类型
func GetResult(i interface{}) string {
	switch i.(type) {
	case string:
		return i.(string)
	case int:
		return strconv.Itoa(i.(int))
	case int64:
		return strconv.FormatInt(i.(int64), 10)
	case float64:
		return strconv.FormatFloat(i.(float64), 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(i.(float64), 'f', -1, 32)
	case []interface{}:
		temp := i.([]interface{})
		var str string
		for val := range temp {
			s := GetResult(temp[val])
			str = str + "," + s
		}
		str = string([]rune(str)[1:])
		str = "[" + str + "]"
		return str
	case interface{}:
		fmt.Println("i come")
		if i!=nil {
			return i.(string)
		}else {
			return ""
		}
	default:
		return ""
	}
}

func GetArr(i interface{}) []string {
	temp := i.([]interface{})
	var str []string
	for val := range temp {
		s := GetResult(temp[val])
		str = append(str, s)
	}
	return str
}
//字符串JSON转map
func PaserStringToMap(msg string) map[string]interface{} {
	var dat map[string]interface{}
	if msg != "" {
		if err := json.Unmarshal([]byte(msg), &dat); err == nil {
			fmt.Println(dat)
		} else {
			fmt.Println(err)
		}
	}
	return dat
}

//字符串JSON数组转map
func PaserStringToMaps(msg string) []map[string]interface{} {
	var dat []map[string]interface{}
	if msg != "" {
		if err := json.Unmarshal([]byte(msg), &dat); err == nil {
		} else {
			fmt.Println(err)
		}
	}
	return dat
}
