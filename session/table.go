package session

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/gofish2020/easyorm/logger"
	"github.com/gofish2020/easyorm/table"
)

// 基于 value 对象，构造表信息
func (s *Session) Model(value interface{}) *Session {
	// 判断表信息 or 类型是否一致
	if s.tableMeta == nil || reflect.TypeOf(value) != reflect.TypeOf(s.tableMeta.Model) {
		meta := table.Parse(value, s.dial)
		if meta == nil {
			logger.Errorf("the value is nil: %v", value)
		} else {
			s.tableMeta = meta
		}
	}
	return s
}

func (s *Session) TableMeta() *table.TableMeta {
	return s.tableMeta
}

func (s *Session) CreateTable() error {
	table := s.TableMeta()

	var columns []string

	for _, field := range table.Fields {
		// 字段名/字段类型/字段约束
		columns = append(columns, fmt.Sprintf("%s %s %s", field.FieldName, field.FieldType, field.FieldTag))
	}

	columnStr := strings.Join(columns, ",")

	_, err := s.Raw(fmt.Sprintf("create table %s (%s)", table.TableName, columnStr)).Exec()
	return err
}

func (s *Session) DropTable() error {

	_, err := s.Raw(fmt.Sprintf("drop table if EXISTS %s", s.TableMeta().TableName)).Exec()
	return err
}

func (s *Session) TableExist() bool {

	sqlStr, sqlParam := s.dial.TableExistSQL(s.TableMeta().TableName)
	row := s.Raw(sqlStr, sqlParam...).QueryRow()

	var res string
	err := row.Scan(&res)
	if err == sql.ErrNoRows {
		return false
	}
	return res == s.TableMeta().TableName
}
