package webgo

import "math"

//反转整数
//例子1: x = 123, return 321
//例子2: x = -123, return -321
func reverse(x int)(num int)  {
	for x!= 0{
		num =num*10+x%10
		x=x/10
	}
	if num > math.MaxInt32 || num<math.MinInt32{
		return 0
	}
	return
}