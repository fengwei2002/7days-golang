package session

import (
	"fmt"
	"reflect"
	"strings"
	"tinygorm/log"

	"tinygorm/schema"
)

/*
用来放置操作数据库表相关的代码
*/

// Model assigns refTable
// 给 refTable 赋值 解析操作是比较耗费时间的，因此将解析的结果保存在成员变量 refTable 中
// 即使 Moder 被调用多次，如果传入的结构体的名称不发生变化，则不会更新 refTable 的值
func (s *Session) Model(value interface{}) *Session {
	// nil or different model, update refTable
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

// RefTable returns a Schema instance that contains all parsed fields
// 返回 refTable 的值， 如果 refTable 未被赋值，这打印错误日志
func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("Model is not set")
	}
	return s.refTable
}

// CreateTable create a table in database with a model
// 实现数据库表的创建，删除和判断是否存在的功能
func (s *Session) CreateTable() error {
	table := s.RefTable()
	var columns []string
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(columns, ",")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", table.Name, desc)).Exec()
	return err
}

// DropTable drops a table with the name of model
func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.RefTable().Name)).Exec()
	return err
}

// HasTable returns true of the table exists
// TableExistSQL 在 dialect 中已经实现了
func (s *Session) HasTable() bool {
	sql, values := s.dialect.TableExistSQL(s.RefTable().Name)
	row := s.Raw(sql, values...).QueryRow()
	var tmp string
	_ = row.Scan(&tmp)
	return tmp == s.RefTable().Name
}
