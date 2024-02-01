package config

import (
	"os"
)

// ValueOrPath is typically used to read Kubernetes secrets (https://kubernetes.io/docs/concepts/configuration/secret/),
// which can be provided to pods using either environment variables or mounted volumes.
// It is an interface to ensure consistency - use NewValueOrPath() function to obtain one
type ValueOrPath interface {
	Value() string
	Path() string
	IsEmpty() bool
	IsFile() bool
}

type valueOrPath struct {
	value string
	path  string
}

// NewValueOrPath returns a new ValueOrPath according to the arguments:
//   - s is either the value itself or a path to a file containing the value (the value is
//     then read as a string)
//   - fileOnly indicates if s must be a path; set it to true for values which
//     should be passed only in files, either for security purposes (k8s secret mounted
//     as a volume), or when the file is needed by whichever library (e.g. a cert file)
//   - shouldRead indicates whether a file should actually be read, so its value is available
//     through this interface (e.g. for a cert file pass false)
func NewValueOrPath(s string, fileOnly, shouldRead bool) (vop ValueOrPath, e error) {
	vopImpl := &valueOrPath{}
	if s != Empty {
		var err error
		if shouldRead {
			var b []byte
			if b, err = os.ReadFile(s); err == nil {
				vopImpl.value = string(b)
			}
		} else {
			if _, err = os.Stat(s); err == nil {
				vopImpl.path = s
			}
		}
		if err != nil {
			if fileOnly {
				e = err
			} else {
				vopImpl.value = s
			}
		}
	}
	vop = vopImpl
	return
}

func (vop *valueOrPath) Value() string {
	return vop.value
}

func (vop *valueOrPath) Path() string {
	return vop.path
}

func (vop *valueOrPath) IsEmpty() bool {
	return vop.value == Empty && !vop.IsFile()
}

func (vop *valueOrPath) IsFile() bool {
	return vop.path != Empty
}
