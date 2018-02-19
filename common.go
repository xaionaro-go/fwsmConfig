package fwsmConfig

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	errNotFound = errors.New("not found")
	errNotImplemented = errors.New("not implemented, yet")
)

func warning(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	fmt.Println("")
}

func removeIndexes(valueI interface{}, idx_ar ...int) interface{} {
	var value reflect.Value

	switch valueT := valueI.(type) {
	case reflect.Value:
		value = valueT
	default:
		value = reflect.ValueOf(valueT)
	}

	idx_count := len(idx_ar)

	if idx_count == 0 {
		return value.Interface()
	}

	idx_ht := make(map[int]bool)
	for _, idx := range idx_ar {
		//fmt.Printf("EXC: %v\n", idx);
		idx_ht[idx] = true
	}

	valueElem := reflect.Indirect(value)

	length := valueElem.Len()
	newLength := length - idx_count

	if newLength < 0 {
		newLength = 0
	}

	var valueElemNew reflect.Value
	switch valueElem.Kind() {
	case reflect.Slice:
		valueElemNew = reflect.MakeSlice(valueElem.Type(), newLength, newLength)
		break
	case reflect.Array:
		arrayType := reflect.ArrayOf(newLength, valueElem.Type())
		valueElemNew = reflect.New(arrayType).Elem()
		break
	default:
		panic(fmt.Errorf("This case is not implemented, yet: %v", valueElem.Kind()))
	}
	i, j := 0, 0
	for i < length {
		for idx_ht[i] {
			i++
		}
		if i >= length {
			break
		}

		valueElemNew.Index(j).Set(valueElem.Index(i))

		i++
		j++
	}

	if !valueElem.CanSet() {
		return valueElemNew.Interface()
	}
	valueElem.Set(valueElemNew)
	return valueElem.Interface()

}
