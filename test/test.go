package _test

import "testing"

func MustPanic(t *testing.T, msg string, f func()) {
	defer func() {
		err := recover()
		if err == nil {
			t.Errorf("function did not panic, wanted %q", msg)
		} else if err != msg {
			t.Errorf("got panic %v, wanted %q", err, msg)
		}
	}()
	f()
}
