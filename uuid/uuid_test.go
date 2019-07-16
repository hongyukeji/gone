package _uuid_test

import (
	"iki-go/uuid"
	"testing"
)

func BenchmarkUUID32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if len(_uuid.UUID32()) != 32 {
			b.Errorf("%s", "生成的uuid长度不为32")
		}
	}
}