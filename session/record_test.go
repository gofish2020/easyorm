package session

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type DataStruct struct {
	Name string
	Id   uint32
}

// &DataStruct{}
func setStructField(dest interface{}) {

	value := reflect.Indirect(reflect.ValueOf(dest))
	var fields []interface{}
	fields = append(fields, value.FieldByName("Name").Addr().Interface()) // 提取结构体字段
	fields = append(fields, value.FieldByName("Id").Addr().Interface())   // 提取结构体字段
	setValue(fields...)
}

func setValue(dests ...interface{}) {

	for _, dest := range dests {
		d := reflect.Indirect(reflect.ValueOf(dest))
		switch d.Kind() {
		case reflect.String:
			d.Set(reflect.ValueOf("hello world"))
		case reflect.Uint32:
			d.Set(reflect.ValueOf(uint32(2023)))
		}
	}
}

func TestSetStruct(t *testing.T) {

	data := &DataStruct{}
	setStructField(data)

	t.Log(data.Name)
	t.Log(data.Id)

	assert.Equal(t, "hello world", data.Name)
	assert.Equal(t, uint32(2023), data.Id)
}

var (
	user1 = &User{18, "Tom1", time.Now()}
	user2 = &User{19, "Tom2", time.Now()}
	user3 = &User{20, "Tom3", time.Now()}
)

func testRecordInit(t *testing.T) *Session {
	t.Helper()
	s := NewSession().Model(&User{})
	err1 := s.DropTable()
	err2 := s.CreateTable()
	_, err3 := s.Insert(user1, user2, user3)
	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatal("failed init test records")
	}
	return s
}
func TestInsert(t *testing.T) {
	testRecordInit(t)
}

func TestFind(t *testing.T) {
	session := testRecordInit(t)

	var users []User
	if err := session.Find(&users); err != nil || len(users) != 3 {
		t.Fatalf("find err:%v", err)
	}
}

func TestFirst(t *testing.T) {
	session := testRecordInit(t)

	var user User
	session.Where("id=?", 30)
	if err := session.First(&user); err != nil && err != sql.ErrNoRows {
		t.Fatalf("first err:%v", err)
	}

	t.Log(user)
}

func TestFind_filter(t *testing.T) {
	session := testRecordInit(t)

	session.Limit(0, 2)
	session.Where("id>=?", 30)
	session.OrderBy("id desc")
	var user []User
	if err := session.Find(&user); err != nil {
		t.Fatalf("find err:%v", err)
	}

	t.Log(user)
}

func TestUpdate(t *testing.T) {
	session := testRecordInit(t)

	res, err := session.Where("id=?", 19).Update("name", "Tooooooom")
	if err != nil {
		t.Fatalf("Update err:%v", err)
	}

	t.Log(res)
}

func TestDeleteAndCount(t *testing.T) {
	session := testRecordInit(t)

	_, err := session.Where("id=?", 18).Delete()
	if err != nil {
		t.Fatalf("Delete err:%v", err)
	}
	count, err := session.Count()
	if err != nil {
		t.Fatalf("Delete err:%v", err)
	}
	assert.Equal(t, 2, count)
	t.Log(count)
}
