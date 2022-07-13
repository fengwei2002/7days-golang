package session

import (
	"database/sql"
	"strings"
	"tinygorm/clause"
	"tinygorm/dialect"
	"tinygorm/log"
	"tinygorm/schema"
)

/*

用于实现和数据库的交互，只实现调用 sql 语句进行原生交互的部分
这部分代码实现在 session/raw.go 中

*/

// Session keep a pointer to sql.DB and provides all execution of all
// kind of database operations.
type Session struct {
	db       *sql.DB         // 使用 sql.Open 方法连接数据库成功之后返回的指针
	sql      strings.Builder // 第二个和第三个变量用来拼接 sql 语句和 sql 语句中占位符的对应 value
	sqlVars  []interface{}   // 用户调用 raw 方法可以改变这两个变量的值
	dialect  dialect.Dialect
	refTable *schema.Schema
	clause   clause.Clause
}

// New creates an instance of Session
func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
	}
}

// Clear initialize the state of a session
// 清空内容但是不关闭会话，让后来的代码可以多次使用
func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = clause.Clause{}
}

// DB returns *sql.DB
func (s *Session) DB() *sql.DB {
	return s.db
}

// Exec raw sql with sqlVars 封装原生方法 exec
// 封装有两个目的
// - 打印统一日志
// - 执行完成之后，清空 (s *Session).sql 和 (s *Session).sqlVars 两个变量
// 这样 session 可以重复使用，开启一次会话，可以执行多次 sql
func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

// QueryRow gets a record from db 封装原生方法 query row
func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

// QueryRows gets a list of records from db 封装原生方法 query rows
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

// Raw appends sql and sqlVars
// 使用 raw 方法修改 sql 语句和 占位符使用的具体 value
func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}
