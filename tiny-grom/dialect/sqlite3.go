package dialect

import (
	"fmt"
	"reflect"
	"time"
)

/*
增加 dialect 对 sqlite 的支持
*/

type sqlite3 struct{}

var _ Dialect = (*sqlite3)(nil)

func init() { // 包在第一次加载的时候，就将 sqlite3 的 dialect 注册到全局中
	RegisterDialect("sqlite3", &sqlite3{})
}

// DataTypeOf Get Data Type for sqlite3 Dialect 将 go 语言的类型映射为 sqlite 的数据类型
func (s *sqlite3) DataTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		return "integer"
	case reflect.Int64, reflect.Uint64:
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "real"
	case reflect.String:
		return "text"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))
}

// TableExistSQL returns SQL that judge whether the table exists in database
// 返回了在 sqlite 中判断表 tableName 是否存在的 sql 语句
func (s *sqlite3) TableExistSQL(tableName string) (string, []interface{}) {
	args := []interface{}{tableName}
	return "SELECT name FROM sqlite_master WHERE type='table' and name = ?", args
}
