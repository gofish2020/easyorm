package table

import (
	"testing"
	"time"

	"github.com/gofish2020/easyorm/dialect"
	"github.com/stretchr/testify/assert"
)

type User struct {
	Id     uint64 `easyorm:"primary key"`
	Name   string
	Sex    bool
	Weight int
	Desc   string
	UpTime time.Time

	password string
}

func (u User) TableName() string {
	return "t_user"
}

func TestParse(t *testing.T) {

	u := &User{}

	sqlite3, _ := dialect.GetDialect("sqlite3")
	meta := Parse(u, sqlite3)

	assert.Equal(t, "t_user", meta.TableName)
	t.Log(meta.TableName)

	assert.Equal(t, 6, len(meta.Fields))
	t.Log(meta.FieldsName)

	assert.Equal(t, "primary key", meta.fieldMap["Id"].FieldTag)
	t.Log(meta.fieldMap["Id"].FieldTag)
}

func TestParseNil(t *testing.T) {

	u := (*User)(nil)

	sqlite3, _ := dialect.GetDialect("sqlite3")
	meta := Parse(u, sqlite3)

	assert.Equal(t, "t_user", meta.TableName)
	t.Log(meta.TableName)

	assert.Equal(t, 6, len(meta.Fields))
	t.Log(meta.FieldsName)

	assert.Equal(t, "primary key", meta.fieldMap["Id"].FieldTag)
	t.Log(meta.fieldMap["Id"].FieldTag)
}
