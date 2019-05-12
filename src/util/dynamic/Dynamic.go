package dynamic

import (
	"reflect"
)

var typeRegistries = make(map[string]reflect.Type)

func RegistierType(elem interface{})  {
	t := reflect.TypeOf(elem).Elem()
	typeRegistries[t.Name()] = t
}

func NewStruct(name string)(interface{}, bool){
	if elem,ok := typeRegistries[name];ok{
		println(elem.Name())
		value := reflect.New(elem)
		println(value.String())
		println(value.Interface())
		return value.Interface(),true
	}
	return nil,false
}