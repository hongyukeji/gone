package opt_test

import (
	"fmt"
	"github.com/wx11055/gone/opt"
	"testing"
)

func TestOptional(t *testing.T) {
	s, err := opt.OfString("123123").String()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(s)
}
