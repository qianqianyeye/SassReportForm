package webgo

import (
	"gopkg.in/go-playground/validator.v8"
	"reflect"
	"time"
	"strconv"
)

func PageValid(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	if num, ok := field.Interface().(int); ok {
		if num<=0 {
			return false
		}
	}
	return true
}

func TimeValid(	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,) bool {
	if _, ok := field.Interface().(time.Time); ok {
		return true
	}
	return false
}

func MerchantIdAndStoreIdValid(v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,) bool {
	if val, ok := field.Interface().(string); ok {
		if val=="all" {
			return true
		}else {
			if i,err:=strconv.Atoi(val);err!=nil{
				return false
			}else if i<=0 {
				return false
			}
		}
	}
	return true
}

