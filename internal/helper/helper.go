package helper

import (
	"bytes"
	"context"
	"errors"
	"reflect"
	"strconv"
)

type TestCase struct {
	Name         string
	Req          interface{}
	ExpectedErr  error
	ExpectedResp interface{}
	Setup        func(ctx context.Context)
}

type SQLFieldMap struct {
	Keys   []string
	Values []interface{}
}

func GeneratePlaceholder(n int) string {
	var b bytes.Buffer
	seperate := ","
	for i := 1; i <= n; i++ {
		if i == n {
			seperate = ""
		}
		b.WriteString("$" + strconv.Itoa(i) + seperate)
	}
	return b.String()
}

func GetSQLFieldMapFromEntity(entity interface{}) (sqlFieldMap SQLFieldMap, err error) {
	value := reflect.ValueOf(entity)
	if value.Kind() == reflect.Ptr && value.Elem().Kind() == reflect.Struct {
		value = value.Elem()
	} else {
		return sqlFieldMap, errors.New("entity must be pointer")
	}
	for i := 0; i < value.Type().NumField(); i++ {
		field := value.Type().Field(i).Tag.Get("sql")
		if field == "" {
			continue
		}
		sqlFieldMap.Keys = append(sqlFieldMap.Keys, field)
		sqlFieldMap.Values = append(sqlFieldMap.Values, value.Field(i).Addr().Interface())
	}
	return sqlFieldMap, nil
}
