package dialect

import "reflect"

/*
使用 dialect 抽象出各个数据库的差异部分
*/

var dialectsMap = map[string]Dialect{}

// Dialect is an interface contains methods that a dialect has to implement
type Dialect interface {
	DataTypeOf(typ reflect.Value) string                    // 用于将 go 语言的数据类型转换为这种数据库的数据类型
	TableExistSQL(tableName string) (string, []interface{}) // 返回某个表中是否存在的 sql 语句 参数是 table name
}

// RegisterDialect register a dialect to the global variable
// 如果增加了对某个数据库的支持，那么调用 RegisterDialect 就可以将他注册到全局的 map 中
func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

// GetDialect Get the dialect from global variable if it exists
func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}
