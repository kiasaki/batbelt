package chk

import "testing"

func Assert(t *testing.T, worked bool) bool {
	if !worked {
		t.Fail()
	}
	return worked
}

func AssertLog(t *testing.T, worked bool, mess string) bool {
	if !worked {
		t.Error(mess)
	}
	return worked
}

func AssertLogf(t *testing.T, worked bool, mess string, args ...interface{}) bool {
	if !worked {
		t.Errorf(mess, args...)
	}
	return worked
}
