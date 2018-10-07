package service

import (
	"github.com/jinzhu/gorm"
	"SaasServiceGo/src/db"
	"SaasServiceGo/src/webgo"
	"database/sql"
	"fmt"
	"strings"
	"reflect"
	"time"
)

func GetDbPageData(pageModel webgo.PageModel) *gorm.DB {
	rd :=db.SqlDB.Offset(pageModel.OffSet).Limit(pageModel.PageSize)
	return rd
}

func GetDbDataMap(rows *sql.Rows)  []map[string]interface{}{
	var m []map[string]interface{}
	cols, _ := rows.Columns()
	for rows.Next(){
		row := make([]interface{}, 0)
		generic := reflect.TypeOf(row).Elem()
		fmt.Println(generic)
		for _ = range cols {
			row = append(row, reflect.New(generic).Interface())
		}
		_ = rows.Scan(row...)
		rowMap := make(map[string]interface{})

		for i, col := range cols {
			switch (*(row[i].(*interface{}))).(type) {
			case []uint8:
				arr := (*(row[i].(*interface{}))).([]uint8)
				rowMap[col] = string(arr)
			case time.Time:
				time := (*(row[i].(*interface{}))).(time.Time)
				rowMap[col] = webgo.GetStringDateTime(time)
			case nil:
				rowMap[col]= ""
			default:
				rowMap[col] = *(row[i].(*interface{}))
			}
		}
		m = append(m, rowMap)
	}
	return m
}

type NullType byte

const (
	_ NullType = iota
	// IsNull the same as `is null`
	IsNull
	// IsNotNull the same as `is not null`
	IsNotNull
)

// sql build where
func whereBuild(where map[string]interface{}) (whereSQL string, vals []interface{}, err error) {
	for k, v := range where {
		ks := strings.Split(k, " ")
		if len(ks) > 2 {
			return "", nil, fmt.Errorf("Error in query condition: %s. ", k)
		}

		if whereSQL != "" {
			whereSQL += " AND "
		}
		strings.Join(ks, ",")
		switch len(ks) {
		case 1:
			//fmt.Println(reflect.TypeOf(v))
			switch v := v.(type) {
			case NullType:
				if v == IsNotNull {
					whereSQL += fmt.Sprint(k, " IS NOT NULL")
				} else {
					whereSQL += fmt.Sprint(k, " IS NULL")
				}
			default:
				whereSQL += fmt.Sprint(k, "=?")
				vals = append(vals, v)
			}
			break
		case 2:
			k = ks[0]
			switch ks[1] {
			case "=":
				whereSQL += fmt.Sprint(k, "=?")
				vals = append(vals, v)
				break
			case ">":
				whereSQL += fmt.Sprint(k, ">?")
				vals = append(vals, v)
				break
			case ">=":
				whereSQL += fmt.Sprint(k, ">=?")
				vals = append(vals, v)
				break
			case "<":
				whereSQL += fmt.Sprint(k, "<?")
				vals = append(vals, v)
				break
			case "<=":
				whereSQL += fmt.Sprint(k, "<=?")
				vals = append(vals, v)
				break
			case "!=":
				whereSQL += fmt.Sprint(k, "!=?")
				vals = append(vals, v)
				break
			case "<>":
				whereSQL += fmt.Sprint(k, "!=?")
				vals = append(vals, v)
				break
			case "in":
				whereSQL += fmt.Sprint(k, " in (?)")
				vals = append(vals, v)
				break
			case "like":
				whereSQL += fmt.Sprint(k, " like ?")
				vals = append(vals, v)
			}
			break
		}
	}
	return
}