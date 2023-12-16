package table

import (
	"reflect"

	"github.com/gofish2020/easyorm/dialect"
)

// 将结构体 转化成 表解雇

// 字段定义 Field
type Field struct {
	FieldName string // 字段名
	FieldType string // 字段类型
	FieldTag  string //
}

// 表定义 TableMeta
type TableMeta struct {
	Model      interface{}
	TableName  string
	Fields     []*Field
	FieldsName []string

	fieldMap map[string]*Field
}

type Tabler interface {
	TableName() string
}

// Parse 将结构体转化成 表/字段/字段类型
func Parse(dest interface{}, d dialect.Dialect) *TableMeta {
	// 是否为nil
	if dest == nil {
		return nil
	}

	// (*User)(nil)
	value := reflect.ValueOf(dest)
	if value.Kind() == reflect.Ptr && value.IsNil() {
		value = reflect.New(value.Type().Elem()) // 取出类型，构造新对象
	}

	meta := &TableMeta{
		Model:    value.Interface(), // 保存对象
		fieldMap: make(map[string]*Field),
	}

	destValue := reflect.Indirect(value)
	destType := destValue.Type()

	// 结构体类型名：就是表名
	m, ok := value.Interface().(Tabler)
	if ok {
		meta.TableName = m.TableName()
	} else {
		meta.TableName = destType.Name()
	}

	for i := 0; i < destType.NumField(); i++ {

		fieldType := destType.Field(i)

		// 非匿名 && 可导出
		if !fieldType.Anonymous && fieldType.IsExported() {
			// 成员变量：就是字段名
			// 成员变量类型：就是字段类型
			field := &Field{
				FieldName: fieldType.Name,
				FieldType: d.ConvertType2DBType(destValue.Field(i)),
				FieldTag:  fieldType.Tag.Get("easyorm"),
			}

			meta.Fields = append(meta.Fields, field)
			meta.FieldsName = append(meta.FieldsName, field.FieldName)
			meta.fieldMap[field.FieldName] = field
		}
	}

	return meta
}

// 提取结构体中的字段值
func (t *TableMeta) ExtractFieldValue(dest interface{}) []interface{} {
	value := reflect.Indirect(reflect.ValueOf(dest))
	var result []interface{}
	for _, fieldName := range t.FieldsName {
		result = append(result, value.FieldByName(fieldName).Interface())
	}
	return result
}
