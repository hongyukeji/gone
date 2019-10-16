package opt

import "fmt"

type OpType int

const (
	Int    OpType = 0
	String OpType = 1
	Float  OpType = 2
)

type Optional struct {
	k string
	v interface{}
	t OpType
}

type Middle struct {
}
type End struct {
}

func OfString(s string) *Optional {
	return &Optional{v: s, t: String}
}
func (o *Optional) String() (string, error) {
	return fmt.Sprintf("%v", o.v), nil
}
