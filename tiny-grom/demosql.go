package tinygorm

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

/*
以下内容是 基础的 sql 语句 和
go 标准库 database/sql 的使用示例
*/

func sqlUseDemo() {
	db, _ := sql.Open("sqlite3", "gee.db")
	// 打开指定的数据表，使用指定的数据库
	defer func() { _ = db.Close() }()
	// 最后要记得将数据库关闭

	_, _ = db.Exec("DROP TABLE IF EXISTS User;")
	_, _ = db.Exec("CREATE TABLE User(Name text);")
	// Exec 用于执行 sql 语句 如果是查询语句 不会返回 相关的记录
	// 所以查询语句通常使用 query 和 queryRow 前者可以返回 多条记录 后者只能返回一条记录

	result, err := db.Exec("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam")
	// 使用 ？ 作为占位符，后面的入参是 sql 语句中的占位符 ？ 对应的值，占位符一般用来预防 sql 注入

	if err == nil {
		affected, _ := result.RowsAffected()
		log.Println(affected)
	}

	row := db.QueryRow("SELECT Name FROM User LIMIT 1")
	// 返回值类型是 *sql.Row row.Scan 接受一个或者多个指针作为参数
	// 可以获取对应列的值

	var name string
	if err := row.Scan(&name); err == nil {
		log.Println(name)
	}
}
