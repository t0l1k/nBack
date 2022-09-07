package ui

import (
	"fmt"
	"image/color"
)

type Preferences map[string]interface{}

func NewPreferences() Preferences {
	return make(Preferences)
}

func (p Preferences) Get(set string) (value interface{}) {
	return p[set]
}

func (p Preferences) Set(set string, value interface{}) {
	p[set] = value
}

func (p Preferences) String() string {
	s := ""
	for k, v := range p {
		s += fmt.Sprintf("%v: %v\n", k, v)
	}
	return s
}

type Theme map[string]color.Color

func NewTheme() Theme {
	return make(Theme)
}

func (t Theme) Get(set string) (value color.Color) {
	return t[set]
}

func (t Theme) Set(set string, value color.Color) {
	t[set] = value
}

func (t Theme) String() string {
	s := ""
	for k, v := range t {
		s += fmt.Sprintf("%v: %v\n", k, v)
	}
	return s
}

type Locale map[string]string

func NewLocale() Locale {
	return make(Locale)
}

func (l Locale) Get(value string) string {
	return l[value]
}

func (l Locale) Set(set, value string) {
	l[set] = value
}

func (l Locale) String() string {
	s := ""
	for k, v := range l {
		s += fmt.Sprintf("%v: %v\n", k, v)
	}
	return s
}
