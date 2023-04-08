package models

import (
	"fmt"
	"net/url"
	"strings"
)

type GetURLRequest struct {
	ID *LinkID `urlparam:"id"`
	// client string
}

func (gr GetURLRequest) VerifyValues() error {
	return VerifyStruct(gr)
}

type LinkID string

func (lid LinkID) String() string {
	return string(lid)
}

func (lid *LinkID) CleanSelf() {
	CleanStringlike(lid)
}

func (lid *LinkID) ValidateSelf() error {
	return nil
}

// ============

type PostURLRequest struct {
	Link        *InputLink   `json:"link"`
	Description *Description `json:"description"`
}

func (pr PostURLRequest) VerifyValues() error {
	return VerifyStruct(pr)
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
