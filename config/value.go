package config

import (
	"os"
)

const (
	Empty = ""
	Comma = ","
	Dot   = "."
)

type ValueOrPath struct {
	value string
	path  string
}

func NewValueOrPath(s string, fileOnly, shouldRead bool) (vop *ValueOrPath, e error) {
	vop = &ValueOrPath{}
	if s != Empty {
		var err error
		if shouldRead {
			var b []byte
			if b, err = os.ReadFile(s); err == nil {
				vop.value = string(b)
			}
		} else {
			if _, err = os.Stat(s); err == nil {
				vop.path = s
			}
		}
		if err != nil {
			if fileOnly {
				e = err
			} else {
				vop.value = s
			}
		}
	}
	return
}

func (vop *ValueOrPath) Value() string {
	return vop.value
}

func (vop *ValueOrPath) Path() string {
	return vop.path
}

func (vop *ValueOrPath) IsEmpty() bool {
	return vop.value == Empty && !vop.IsFile()
}

func (vop *ValueOrPath) IsFile() bool {
	return vop.path != Empty
}
