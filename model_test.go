package model

import (
	"errors"
	"fmt"
	"testing"
)

type MyString string

////////////////////////////////////////////////////////////////////////////////////////////////////
type Human struct {
	Model
	Name MyString `model:"name"`
	Age  int      `model:"age"`
}

func (this *Human) CleanedName(n string) (MyString, error) {
	if len(n) > 0 {
		return MyString(n), nil
	}
	return "", errors.New("随便给点吧")
}

////////////////////////////////////////////////////////////////////////////////////////////////////
type Class struct {
	ClassName string `model:"class_name"`
}

////////////////////////////////////////////////////////////////////////////////////////////////////
type Student struct {
	Human
	Number int `model:"number"`
	Class  Class
}

var source = map[string]interface{}{"name": "Yangfeng", "age": 123, "number": 1234, "class_name1": "adfsf"}

func TestBindPoint(t *testing.T) {
	var s *Student
	fmt.Println(Bind(source, &s))
	if s != nil {
		fmt.Println(s.Name, s.Age, s.Number)
	}
}
