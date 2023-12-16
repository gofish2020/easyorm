package session

import "testing"

type HookStruct struct {
	Password string
}

// 就是对象的回调函数
func (h *HookStruct) AfterQuery(s *Session) (interface{}, error) {
	h.Password = "******"
	return nil, nil
}

var hook1 = &HookStruct{
	Password: "123",
}

var hook2 = &HookStruct{
	Password: "123",
}
var hook3 = &HookStruct{
	Password: "123",
}

func testHookRecordInit(t *testing.T) *Session {
	t.Helper()
	s := NewSession().Model(&HookStruct{})
	err1 := s.DropTable()
	err2 := s.CreateTable()
	_, err3 := s.Insert(hook1, hook2, hook3)
	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatal("failed init test records")
	}
	return s
}

func TestHook(t *testing.T) {
	s := testHookRecordInit(t)

	var hooks []HookStruct
	err := s.Find(&hooks)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hooks)
}
