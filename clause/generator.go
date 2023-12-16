package clause

import (
	"fmt"
	"strings"
)

type generator func(values ...interface{}) (string, []interface{})

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderBy
	generators[UPDATE] = _update
	generators[DELETE] = _delete
	generators[COUNT] = _count
}

// insert into #table ( #fieldname )
func _insert(values ...interface{}) (string, []interface{}) {

	tableName := values[0].(string)

	fields := strings.Join(values[1].([]string), ",")

	return fmt.Sprintf("insert into %s (%s)", tableName, fields), nil
}

// values (#value"),(#value)
func _values(values ...interface{}) (string, []interface{}) {

	var sql strings.Builder
	sql.WriteString(" values ")

	placeholderStr := ""

	var sqlParam []interface{}

	for i, value := range values {

		val := value.([]interface{}) // 每个是个切片

		if placeholderStr == "" {
			placeholderStr = placeholder(len(val))
		}

		sql.WriteString(fmt.Sprintf("(%s)", placeholderStr)) // 构造 (?,?,?)

		if i != len(values)-1 {
			sql.WriteString(",")
		}

		sqlParam = append(sqlParam, val...)
	}
	return sql.String(), sqlParam
}

func placeholder(num int) string {

	var result []string
	for i := 0; i < num; i++ {
		result = append(result, "?")
	}
	return strings.Join(result, ",")
}

// tablename  []string{}
func _select(values ...interface{}) (string, []interface{}) {

	tableName := values[0].(string)

	return fmt.Sprintf(" select  %s from %s ",
		strings.Join(values[1].([]string), ","), tableName), nil
}

func _limit(values ...interface{}) (string, []interface{}) {
	return " limit ?,? ", []interface{}{values[0], values[1]}
}

// order by id asc,name desc
func _orderBy(values ...interface{}) (string, []interface{}) {

	var val []string

	for _, value := range values {
		val = append(val, value.(string))
	}

	return fmt.Sprintf(" order by %s ", strings.Join(val, ",")), nil
}

// where id >? and name = ?
func _where(values ...interface{}) (string, []interface{}) {

	format := values[0].(string)

	return fmt.Sprintf(" where %s ", format), values[1:]
}

// update $table set name=?,id=?
func _update(values ...interface{}) (string, []interface{}) {

	tableName := values[0].(string)
	m := values[1].(map[string]interface{})

	var (
		fieldName  []string
		fieldValue []interface{}
	)
	for k, v := range m {

		fieldName = append(fieldName, k+"=?")
		fieldValue = append(fieldValue, v)
	}

	return fmt.Sprintf(" update %s set %s ", tableName, strings.Join(fieldName, ",")), fieldValue
}

// _delete from $table

func _delete(values ...interface{}) (string, []interface{}) {

	tableName := values[0].(string)
	return fmt.Sprintf(" delete from %s ", tableName), nil
}

// _count

func _count(values ...interface{}) (string, []interface{}) {

	return _select(values[0], []string{" count(1) as count "})
}
