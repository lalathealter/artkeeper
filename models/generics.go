package models

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"

	"github.com/lalathealter/artkeeper/config"
)

type Stringlike interface {
	string |
		InputLink | Description |
		ResourceID | UserID |
		StringifiedInt
	String() string
}

type Cleanable interface {
	CleanSelf()
}

type Validatable interface {
	ValidateSelf() error
}

type Message interface {
	VerifyValues() error
	Call(*sql.DB) (DBResult, error)
}

type DBResult interface{}

func ReflectCastedStringlike(payload string, reference interface{}) (reflect.Value, error) {
	switch (reference).(type) {
	case *ResourceID:
		lid := ResourceID(payload)
		return reflect.ValueOf(&lid), nil
	case *StringifiedInt:
		num := StringifiedInt(payload)
		return reflect.ValueOf(&num), nil
	default:

		return reflect.Value{}, fmt.Errorf("casting to Stringlike isn't implemented in the type switch with %T;", reference)
	}
}

func CleanStringlike[T Stringlike](fieldptr *T) {
	cleanedval := config.Sanitbypolicy(string(*fieldptr))
	*fieldptr = T(cleanedval)
}

func VerifyFieldValue(val reflect.Value) error {
	actualval := reflect.Indirect(val)
	if !actualval.IsValid() {
		return fmt.Errorf("not enough values for the struct's fields provided;")
	}

	currtype := actualval.Type()
	fmt.Printf("IN: '%v' of type %v\n", actualval, currtype)
	c, ok := val.Interface().(Cleanable)
	if !ok {
		return (fmt.Errorf("%v - cleaning method is not implemented", currtype))
	}
	c.CleanSelf()

	v, ok := val.Interface().(Validatable)
	if !ok {
		return (fmt.Errorf("%v - validating method is not implemented", currtype))
	}
	err := v.ValidateSelf()
	if err != nil {
		return err
	}

	return nil
}

func VerifyStruct[T Message](vstruct T) error {

	values := reflect.ValueOf(vstruct)

	for i := 0; i < reflect.ValueOf(vstruct).NumField(); i++ {
		field := values.Field(i)

		if field.Kind() == reflect.Slice || field.Kind() == reflect.Array {
			for i := 0; i < field.Len(); i++ {
				err := VerifyFieldValue(field.Index(i))
				fmt.Println(field.Index(i), err)
				if err != nil {
					return err
				}
			}
		} else {
			err := VerifyFieldValue(field)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func isValidInt[T Stringlike](in T) error {
	_, err := strconv.Atoi(in.String())
	return err
}


func ExtractFieldPointers[T any](in *T) []any {
	iter := reflect.ValueOf(in).Elem()
	fieldptrs := make([]any, iter.NumField())
	for i := 0; i < iter.NumField(); i++ {
		fieldptrs[i] = iter.Field(i).Addr().Interface()
	}
	return fieldptrs
}

// func ExtractFieldValues[T any](in *T) []any {
// 	iter := reflect.ValueOf(in).Elem()
// 	fieldvals := make([]any, iter.NumField())
// 	for i := 0; i < iter.NumField(); i++ {
// 		fieldvals[i] = iter.Field(i).Interface()
// 	}
// 	return fieldvals
// }
