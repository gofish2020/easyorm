package easyorm

import (
	"database/sql"

	"github.com/gofish2020/easyorm/dialect"
	"github.com/gofish2020/easyorm/logger"
	"github.com/gofish2020/easyorm/session"
	_ "github.com/mattn/go-sqlite3"
)

type Engine struct {
	db   *sql.DB
	dial dialect.Dialect
}

func NewEngine(driverName, dataSourceName string) (*Engine, error) {

	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	dial, ok := dialect.GetDialect(driverName)
	if !ok {
		logger.Errorf("dialect %s Not Found", driverName)
		return nil, nil
	}
	e := &Engine{db: db, dial: dial}
	logger.Infof("connect %s success", driverName)
	return e, nil
}

func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		logger.Errorf("db close err:%v", err)
		return
	}
	logger.Infof("db close success")
}

// NewSession creates a new session for next operations
func (engine *Engine) NewSession() *session.Session {
	return session.New(engine.db, engine.dial)
}

// TxFunc will be called between tx.Begin() and tx.Commit()
type TxFunc func(*session.Session) (interface{}, error)

// Transaction executes sql wrapped in a transaction, then automatically commit if no error occurs
func (engine *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := engine.NewSession()
	if err := s.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			_ = s.Rollback() // err is non-nil; don't change it
		} else {
			err = s.Commit() // err is nil; if Commit returns error update err
		}
	}()

	return f(s)
}
