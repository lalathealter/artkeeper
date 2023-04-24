package models

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type StringifiedInt string

func (i *StringifiedInt) CleanSelf() {
	CleanStringlike(i) // may be unnecessary (??)
	v, _ := strconv.Atoi(string(*i))
	(*i) = StringifiedInt(strconv.Itoa(v))
}

func (i *StringifiedInt) ValidateSelf() error {
	return nil
}

func (i StringifiedInt) String() string {
	return string(i)
}

func (gr GetURLRequest) VerifyValues() error {
	return VerifyStruct(gr)
}

type ResourceID string

func (lid ResourceID) String() string {
	return string(lid)
}

func (lid *ResourceID) CleanSelf() {
	CleanStringlike(lid)
}

func (lid *ResourceID) ValidateSelf() error {
	if *lid == "" {
		return nil // allowing empty values for get all method
	}
	return isValidInt(*lid)
}

type InputLink string

var validschemes = map[string]bool{
	"http":  true,
	"https": true,
}

func (il InputLink) ValidateSelf() error {
	link := string(il)
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

func (il *InputLink) CleanSelf() {
	CleanStringlike(il)
}

func (il InputLink) String() string {
	return string(il)
}

type Description string

func (d Description) ValidateSelf() error {
	return nil
}

func (d *Description) CleanSelf() {
	CleanStringlike(d)
}

func (d Description) String() string {
	return string(d)
}

type UserID string

func (uid UserID) ValidateSelf() error {
	// authenticate + authorize here?
	return isValidInt(uid)
}

func (uid *UserID) CleanSelf() {
	CleanStringlike(uid)
}

func (uid UserID) String() string {
	return string(uid)
}
