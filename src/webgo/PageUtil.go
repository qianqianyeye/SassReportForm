package webgo

import (
	"strconv"
)
type PageModel struct {
	PageSize int `json:"page_size"`
	Current int `json:"current"`
	Total int `json:"total"`
	OffSet string  `json:"-"`
	Status []string `json:"-"`
}

func GetPageInfo(pageSize string,current string) PageModel {
	if pageSize =="" {
		pageSize="10"
	}else if current=="" {
		current="1"
	}
	ps ,_:=strconv.Atoi(pageSize)
	cr ,_:=strconv.Atoi(current)
	off :=(cr-1)*ps
	offset:=strconv.Itoa(off)
	var PageModel PageModel
	PageModel.PageSize=ps
	PageModel.Current=cr
	PageModel.OffSet=offset
	return PageModel
}
