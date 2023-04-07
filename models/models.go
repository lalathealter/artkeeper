package models

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/lalathealter/artkeeper/config"
)

// ====== INTERFACES =====

type Stringlike interface {
	Userlink | Description | string
	String() string
}

func Cleanstringlike[T Stringlike](val *T) {
	*val = T(config.Sanitbypolicy(string(*val)))
}

type Cleanable interface {
	Cleanself()
}
type Validatable interface {
	Validateself() error
}

type Message interface {
	Verifyvalues() error
}

// ====== STRUCTS ========

type Posturl struct {
	Link        *Userlink    `json:"link"`
	Description *Description `json:"description"`
}

func (p Posturl) Verifyvalues() error {

	values := reflect.ValueOf(p)

	for i := 0; i < reflect.ValueOf(p).NumField(); i++ {
		val := values.Field(i)
		actval := reflect.Indirect(val)

		fmt.Println("IN", actval)
		c, ok := val.Interface().(Cleanable)
		if !ok {
			return (fmt.Errorf("%v %v - cleaning method is not implemented", actval.Type(), actval))
		}
		c.Cleanself()

		v, ok := val.Interface().(Validatable)
		if !ok {
			return (fmt.Errorf("%v - validating method is not implemented", actval.Type()))
		}
		err := v.Validateself()
		if err != nil {
			return err
		}
	}
	return nil
}

// ======= FIELDS ========

type Userlink string

var validschemes = map[string]bool{
	"http":  true,
	"https": true,
}

func (l Userlink) Validateself() error {
	link := string(l)
	if _, err := url.ParseRequestURI(link); err != nil {
		return err
	}

	u, err := url.Parse(link)
	if err != nil {
		return err
	}

	if !validschemes[u.Scheme] {
		return (fmt.Errorf("parse %v: URL's protocol scheme is unacceptable", link))
	}

	splitports := strings.Split(u.Host, ":")
	hname := splitports[0]

	domains := len(strings.Split(hname, "."))
	ports := len(splitports) - 1
	hasoffdots := (hname[0] == '.' || hname[len(hname)-1] == '.')
	if ports > 1 || domains <= 1 || hasoffdots {
		return (fmt.Errorf("parse %v: URL has incorrect hostname", link))
	}

	return nil
}

func (l *Userlink) Cleanself() {
	Cleanstringlike(l)
}

func (l Userlink) String() string {
	return string(l)
}

type Description string

func (l Description) Validateself() error {
	return nil
}

func (l *Description) Cleanself() {
	Cleanstringlike(l)
}

func (l Description) String() string {
	return string(l)
}
