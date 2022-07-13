package tinygorm

import (
	"database/sql"
	"log"
)

/*
标准 sql 库使用 事务的方式举例
*/

func Transaction() {
	db, _ := sql.Open("sqlite3", "koo.db")
	defer func() { _ = db.Close() }()
	_, _ = db.Exec("CREATE TABLE IF NOT EXISTS User(`Name` text);")

	tx, _ := db.Begin() // 调用 db.Begin() 得到 *sql.Tx 对象，使用 tx.Exec() 执行一系列操作
	_, err1 := tx.Exec("INSERT INTO User(`Name`) VALUES (?)", "Tom")
	_, err2 := tx.Exec("INSERT INTO User(`Name`) VALUES (?)", "Jack")
	if err1 != nil || err2 != nil {
		_ = tx.Rollback() // 如果发生错误，使用 RollBack 回滚，
		log.Println("Rollback", err1, err2)
	} else {
		_ = tx.Commit() // 如果没有发生错误，则通过 tx.Commit() 提交
		log.Println("Commit")
	}
}
