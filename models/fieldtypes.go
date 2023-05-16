package models

import (
	"encoding/hex"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
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

type ResourceID string

func (lid ResourceID) String() string {
	return string(lid)
}

func (lid *ResourceID) CleanSelf() {
	CleanStringlike(lid)
}

func (lid *ResourceID) ValidateSelf() error {
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

type Tag string 

const MAX_TAG_LEN = 64
func (t Tag) ValidateSelf() error {
	if len(t) > MAX_TAG_LEN {
		return fmt.Errorf("tag name is too long")
	}
	return nil
}

func (t *Tag) CleanSelf() {
	CleanStringlike(t)
}

func (t Tag) String() string {
	return string(t)
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

type Username string 

const MAX_NAME_LEN = 36
func (un Username) ValidateSelf() error {
	if len(un) > MAX_NAME_LEN {
		return fmt.Errorf("username is too long")
	}
	return nil
}

func (un *Username) CleanSelf() {
	CleanStringlike(un)
}

func (un Username) String() string {
	return string(un)
}


type Password string 

const MIN_PASS_LEN = 8
func (pass Password) ValidateSelf() error {
	if len(pass)< MIN_PASS_LEN {
		return fmt.Errorf("password is too short")
	}
	return nil
}

func (pass *Password) CleanSelf() {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*pass), bcrypt.DefaultCost)
	if err != nil {
		log.Panicln(err)
	}
	nextPass := Password(hex.EncodeToString(hashedPassword))
	*pass = nextPass
}

func (pass Password) String() string {
	return string(pass)
}
