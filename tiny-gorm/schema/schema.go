package schema

import (
	"go/ast"
	"reflect"
	"tinygorm/dialect"
)

/*
实现最主要的 object to table
*/

// Field represents a column of database
type Field struct {
	Name string // 字段名
	Type string // 类型 type
	Tag  string // 约束条件 tag 例如不能为空，是主键等等
}

// 例如：CREATE TABLE `User` (`Name` text PRIMARY KEY, `Age` integer);
// type User struct {
// 	Name string `tinyorm:"PRIMARY KEY"`
// 	Age int
// }

// Schema represents a table of database
type Schema struct {
	Model      interface{}       // 被映射的对象
	Name       string            // 表名 name
	Fields     []*Field          // 字段
	FieldNames []string          // 包含所有的字段名，列名
	fieldMap   map[string]*Field // 记录字段名和 field 的映射关系，方便之后直接使用，无需遍历 fields
}

// GetField returns field by name
func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

type ITableName interface {
	TableName() string
}

// Parse a struct to a Schema instance
// 将任意的对象解析为 schema 实例
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	// TypeOf() 和 valueOf() 是 reflect 的最基本的两个方法
	// 用来返回入参的类型和值

	var tableName string
	t, ok := dest.(ITableName)
	if !ok {
		tableName = modelType.Name()
	} else {
		tableName = t.TableName()
	}
	schema := &Schema{
		Model:    dest,
		Name:     tableName,
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
				// 使用 reflect.Indirect 来获取指针指向的实例
			}
			if v, ok := p.Tag.Lookup("tinygorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}

// RecordValues return the values of dest 's member variables
//  根据数据库中列的值，从对象中找出对应的值，按照顺序平铺
/*
INSERT 对应的 SQL 语句一般是这样的：
INSERT INTO table_name(col1, col2, col3, ...) VALUES
    (A1, A2, A3, ...),
    (B1, B2, B3, ...),
    ...
在 ORM 框架中期望 Insert 的调用方式如下：
s := orm.NewEngine("sqlite3", "koo.db").NewSession()
u1 := &User{Name: "Tom", Age: 18}
u2 := &User{Name: "Sam", Age: 25}
s.Insert(u1, u2, ...)
*/
func (schema *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range schema.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}
