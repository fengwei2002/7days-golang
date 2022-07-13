package clause

import (
	"reflect"
	"testing"
)

func TestClause_Set(t *testing.T) {
	// 测试一下将结构体映射为 insert 语句 的函数能不能行
	var clause Clause
	clause.Set(INSERT, "User", []string{"Name", "Age"})
	sql := clause.sql[INSERT]
	vars := clause.sqlVars[INSERT]
	t.Log(sql, vars)
	if sql != "INSERT INTO User (Name,Age)" || len(vars) != 0 {
		t.Fatal("failed to get clause")
	}
}

func testSelect(t *testing.T) {
	// 将 clause 中的内容映射为 select 语句 带有后缀选项
	var clause Clause
	clause.Set(LIMIT, 3)
	clause.Set(SELECT, "User", []string{"*"})
	clause.Set(WHERE, "Name = ?", "Tom")
	clause.Set(ORDERBY, "Age ASC")
	sql, vars := clause.Build(SELECT, WHERE, ORDERBY, LIMIT)
	t.Log(sql, vars)
	if sql != "SELECT * FROM User WHERE Name = ? ORDER BY Age ASC LIMIT ?" {
		t.Fatal("failed to build SQL")
	}
	if !reflect.DeepEqual(vars, []interface{}{"Tom", 3}) {
		t.Fatal("failed to build SQLVars")
	}
}

func TestClause_Build(t *testing.T) {
	t.Run("select", func(t *testing.T) {
		testSelect(t)
	})
}
