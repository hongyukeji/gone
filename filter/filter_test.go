package filter_test

import (
	"github.com/dxvgef/filter"
	"log"
	"testing"
)

func TestFilterValue(t *testing.T) {
	password, err := filter.FromString("123567", "密码").Trim().MinLength(6).MaxLength(32).String()
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(password)
}
func TestFilterSetValue(t *testing.T) {
	var password string
	err := filter.Set(
		&password, filter.FromString("Abc123-", "密码").
			MinLength(6).MaxLength(32).HasLetter().HasUpper().HasDigit().HasSymbol(),
	)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func TestMultiFilter(t *testing.T) {
	var reqData struct {
		password string
		age      int16
	}
	err := filter.MSet(
		filter.El(&reqData.password,
			filter.FromString("Abc123-", "密码").
				Required().MinLength(6).MaxLength(32).HasLetter().HasUpper().HasDigit().HasSymbol(),
		),
		filter.El(&reqData.age,
			filter.FromString("3", "年龄").
				IsDigit().MinInteger(18)),
	)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("密码", reqData.password)
	log.Println("年龄", reqData.age)
}
func TestSilent(t *testing.T) {
	var value string
	filter.Set(&value,
		filter.FromString("12341", "数字").MinLength(10, "长度不足10位"),
	)
	t.Log(value)
}
