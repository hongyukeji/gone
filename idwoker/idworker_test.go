package idworker_test

import (
	"gone/idwoker"
	"strconv"
	"testing"
)

// todo 很奇怪连续生成4097个id都是19位当生成4098个就成了18位
func TestIdWorker_InitIdWorker(t *testing.T) {
	currWorker := &idworker.IdWorker{}
	currWorker.InitIdWorker(1000, 1)
	for i := 0; i < 4098; i++ {
		id, err := currWorker.NextId()
		if err != nil {
			t.Error(err)
		}
		//t.Log(id)
		if len([]byte(strconv.Itoa(int(id)))) != 19 {
			t.Error("长度不是19位 ==>>", id, "||", strconv.Itoa(int(id)))
			continue
		}
	}
}
