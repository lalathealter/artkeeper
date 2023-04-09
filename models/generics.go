package models

import (
	"fmt"
	"reflect"

	"github.com/lalathealter/artkeeper/config"
)

// ====== INTERFACES =====

type Stringlike interface {
	InputLink | Description | LinkID | UserID | string
	String() string
}

func ReflectCastedStringlike(payload string, reference interface{}) (reflect.Value, error) {
	switch (reference).(type) {
	case *LinkID:
		lid := LinkID(payload)
		return reflect.ValueOf(&lid), nil
	default:
		return reflect.Value{}, fmt.Errorf("casting to Stringlike isn't implemented in the type switch with %T;", reference)
	}
}

func CleanStringlike[T Stringlike](fieldptr *T) {
	cleanedval := config.Sanitbypolicy(string(*fieldptr))
	*fieldptr = T(cleanedval)
}

func VerifyStruct[T Message](vstruct T) error {

	values := reflect.ValueOf(vstruct)

	for i := 0; i < reflect.ValueOf(vstruct).NumField(); i++ {
		val := values.Field(i)
		actualval := reflect.Indirect(val)
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
	}
	return nil
}

type Cleanable interface {
	CleanSelf()
}
type Validatable interface {
	ValidateSelf() error
}

type Message interface {
	VerifyValues() error
}

// ============
