package session

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/gofish2020/easyorm/dialect"
	_ "github.com/mattn/go-sqlite3"
)

var (
	TestDB      *sql.DB
	TestDial, _ = dialect.GetDialect("sqlite3")
)

func TestMain(m *testing.M) {
	TestDB, _ = sql.Open("sqlite3", "../gee.db")
	code := m.Run()
	_ = TestDB.Close()
	os.Exit(code)
}

func NewSession() *Session {
	return New(TestDB, TestDial)
}

type User struct {
	Id uint32 `easyorm:"primary key"`

	Name   string
	UpTime time.Time
}

func TestCreataTable(t *testing.T) {
	s := NewSession()
	s.Model(&User{})

	err := s.DropTable()
	t.Log(err)
	_ = s.CreateTable()
	if !s.TableExist() {
		t.Fatal("Failed to create table User")
	}
}
