package session

import (
	"errors"
	"reflect"

	"github.com/gofish2020/easyorm/clause"
)

// s.Insert(&User{},&User{})
func (s *Session) Insert(values ...interface{}) (int64, error) {
	if len(values) == 0 {
		return 0, errors.New("param is empty")
	}
	// init table
	s.Model(values[0])
	// clause insert
	s.clause.Set(clause.INSERT, s.TableMeta().TableName, s.TableMeta().FieldsName)

	valueSlice := []interface{}{}
	for _, value := range values {
		// 每个 []interface{} 作为 interface{}
		s.CallMethod(BeforeInsert, value)
		valueSlice = append(valueSlice, s.TableMeta().ExtractFieldValue(value))
	}
	// clause values
	s.clause.Set(clause.VALUES, valueSlice...)
	sql, sqlParam := s.clause.Build(clause.INSERT, clause.VALUES)
	// exec sql
	result, err := s.Raw(sql, sqlParam...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterInsert, nil)
	return result.RowsAffected()
}

// s.Find(&[]User{})  将查询的结果 保存到 切片values中
func (s *Session) Find(values interface{}) error {
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destItemType := destSlice.Type().Elem()

	tableMeta := s.Model(reflect.New(destItemType).Interface()).TableMeta()

	s.clause.Set(clause.SELECT, tableMeta.TableName, tableMeta.FieldsName)

	sql, sqlParam := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)

	rows, err := s.Raw(sql, sqlParam...).Query()
	if err != nil {
		return err
	}

	for rows.Next() {

		record := reflect.New(destItemType).Elem()

		// 遍历结构体中存放结果的字段
		var rowRecord []interface{}
		for _, fieldName := range tableMeta.FieldsName {
			rowRecord = append(rowRecord, record.FieldByName(fieldName).Addr().Interface())
		}
		// 将结构保存到字段中（等价于保存到结构体中）
		if err = rows.Scan(rowRecord...); err != nil {
			return err
		}

		s.CallMethod("AfterQuery", record.Addr().Interface())

		destSlice.Set(reflect.Append(destSlice, record))
	}
	return rows.Close()
}

// s.First(&User{})
func (s *Session) First(value interface{}) error {

	destValue := reflect.Indirect(reflect.ValueOf(value))

	tableMeta := s.Model(value).TableMeta()

	s.Limit(0, 1)
	s.clause.Set(clause.SELECT, tableMeta.TableName, tableMeta.FieldsName)

	sql, sqlParam := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	row := s.Raw(sql, sqlParam...).QueryRow()

	var rowRecord []interface{}
	for _, fieldName := range tableMeta.FieldsName {
		rowRecord = append(rowRecord, destValue.FieldByName(fieldName).Addr().Interface())
	}

	err := row.Scan(rowRecord...)
	if err != nil {
		return err
	}

	s.CallMethod(AfterQuery, value)
	return nil
}

// s.Limit(0,1)
func (s *Session) Limit(offset, num int) *Session {
	s.clause.Set(clause.LIMIT, offset, num)
	return s
}

// s.OrderBy("id asc","name desc")
func (s *Session) OrderBy(values ...interface{}) *Session {
	s.clause.Set(clause.ORDERBY, values...)
	return s
}

// s.Where("Id > ? and name = ?",30,"tom")
func (s *Session) Where(format string, values ...interface{}) *Session {

	var vars []interface{}

	vars = append(vars, format)
	vars = append(vars, values...)
	s.clause.Set(clause.WHERE, vars...)
	return s
}

// s.Update("Name","Tom4","Id",40)
func (s *Session) Update(values ...interface{}) (int64, error) {

	m := make(map[string]interface{})
	num := len(values)
	if num&(num-1) == 0 { // 2的整数倍
		for i := 0; i < num; i = i + 2 {
			m[values[i].(string)] = values[i+1]
		}
	} else {
		return 0, errors.New("param invaild")
	}

	s.clause.Set(clause.UPDATE, s.TableMeta().TableName, m)
	sql, sqlParam := s.clause.Build(clause.UPDATE, clause.WHERE)
	result, err := s.Raw(sql, sqlParam...).Exec()
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// s.Delete()
func (s *Session) Delete() (int64, error) {

	s.clause.Set(clause.DELETE, s.TableMeta().TableName)
	sql, sqlParam := s.clause.Build(clause.DELETE, clause.WHERE)
	result, err := s.Raw(sql, sqlParam...).Exec()
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (s *Session) Count() (int64, error) {

	s.clause.Set(clause.COUNT, s.TableMeta().TableName)

	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vars...).QueryRow()
	var tmp int64
	if err := row.Scan(&tmp); err != nil {
		return 0, err
	}
	return tmp, nil
}
