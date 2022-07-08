package tinygorm

import (
	"database/sql"
	"tinygorm/dialect"

	"tinygorm/log"
	"tinygorm/session"
)

/*

session 负责和数据库的交互

那么交互之前的准备工作，比如连接/测试数据库
交互之后的收尾工作 （关闭连接）等就交给 engine 实现

engine 是 tinygorm 和 用户交互的入口

*/

// Engine is the main struct of geeorm, manages all db sessions and transactions.
type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}
	// Send a ping to make sure the database connection is alive.
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}
	// make sure the specific dialect exists
	dial, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s Not Found", driver)
		return
	}
	e = &Engine{db: db, dialect: dial}
	log.Info("Connect database success")
	return
}

func (engine *Engine) NewSession() *session.Session {
	return session.New(engine.db, engine.dialect)
}
