package test

import (
	"fmt"
	"github.com/jimersylee/iris-seed/commons/strcase"
	"testing"
	"time"
)

func TestDate2unix(t *testing.T) {
	fmt.Println(time.Now().Unix())
}

func TestStrcase(t *testing.T) {
	s := "Any kind_of str"
	fmt.Println(strcase.ToLowerCamel(s))
}
