package session

import (
	"database/sql"
	"strings"

	"github.com/gofish2020/easyorm/clause"
	"github.com/gofish2020/easyorm/dialect"
	"github.com/gofish2020/easyorm/logger"
	"github.com/gofish2020/easyorm/table"
)

type Session struct {
	db *sql.DB
	tx *sql.Tx

	// table
	tableMeta *table.TableMeta

	dial dialect.Dialect
	// clause

	clause clause.Clause

	// sql query
	sql strings.Builder
	// sql param
	sqlParam []interface{}
}

func New(db *sql.DB, dial dialect.Dialect) *Session {
	session := Session{db: db, dial: dial, clause: clause.Clause{}}
	return &session
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlParam = nil
	s.clause = clause.Clause{}
}

// purpose: for s.tx/s.db
type BaseDB interface {
	Query(query string, args ...any) (*sql.Rows, error)
	Exec(query string, args ...any) (sql.Result, error)
	QueryRow(query string, args ...any) *sql.Row
}

var _ BaseDB = (*sql.DB)(nil)
var _ BaseDB = (*sql.Tx)(nil)

func (s *Session) DB() BaseDB {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

// 实际执行 sql 和 sql参数
func (s *Session) Exec() (sql.Result, error) {
	defer s.Clear()
	logger.Info(s.sql.String(), s.sqlParam)

	return s.DB().Exec(s.sql.String(), s.sqlParam...)
}

func (s *Session) Query() (*sql.Rows, error) {
	defer s.Clear()
	logger.Info(s.sql.String(), s.sqlParam)
	return s.DB().Query(s.sql.String(), s.sqlParam...)
}

func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	logger.Info(s.sql.String(), s.sqlParam)
	return s.DB().QueryRow(s.sql.String(), s.sqlParam...)
}

// 汇总 sql 和 sql参数
func (s *Session) Raw(sql string, sqlParam ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString("")
	s.sqlParam = append(s.sqlParam, sqlParam...)
	return s
}
