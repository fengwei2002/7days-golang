package main

import (
	_ "github.com/mattn/go-sqlite3"
	"tinygorm"
)

func main() {
	tinygorm.SqlUseDemo()
	//
	// engine, _ := tinygorm.NewEngine("sqlite3", "koo.db")
	// defer engine.Close()
	// s := engine.NewSession()
	// _, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	// _, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	// _, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	// result, _ := s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	// count, _ := result.RowsAffected()
	// fmt.Printf("Exec success, %d affected\n", count)
}
