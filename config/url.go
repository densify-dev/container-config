package config

import (
	"fmt"
	"net/url"
	"strings"
)

type UrlConfig struct {
	Scheme            string `yaml:"scheme"`
	Host              string `yaml:"host"`
	Port              uint64 `yaml:"port"`
	Username          string `yaml:"username,omitempty"`
	Password          string `yaml:"password,omitempty"`
	EncryptedPassword string `yaml:"encrypted_password,omitempty"`
	Url               string `yaml:"-"`
}

const (
	Http                  = "http"
	Https                 = Http + "s"
	NoPort         uint64 = 99999 // 0 is a valid (system) port, so need something > maxPort
	maxPort        uint64 = 65535
	hostPortFormat        = "%s:%d"
)

var validSchemes = map[string]bool{Http: true, Https: true}

func (uc *UrlConfig) numMandatory() (n int) {
	if uc != nil {
		if uc.Scheme != Empty {
			n++
		}
		if uc.Host != Empty {
			n++
		}
	}
	return
}

func (uc *UrlConfig) finalize() (err error) {
	switch uc.numMandatory() {
	case 0:
		return
	case 1:
		err = fmt.Errorf("invalid UrlConfig")
		return
	}
	var sc string
	if sc, err = validScheme(uc.Scheme); err != nil {
		return
	}
	var h string
	if uc.Port == NoPort {
		h = uc.Host
	} else {
		if err = validatePort(uc.Port); err == nil {
			h = fmt.Sprintf(hostPortFormat, uc.Host, uc.Port)
		} else {
			return
		}
	}
	u := &url.URL{
		Scheme: sc,
		Host:   h,
	}
	uc.Url = u.String()
	return
}

func validScheme(scheme string) (s string, err error) {
	s = strings.ToLower(scheme)
	if ok := validSchemes[s]; !ok {
		err = fmt.Errorf("invalid scheme: %s", scheme)
	}
	return
}

func validatePort(port uint64) (err error) {
	if port > maxPort {
		err = fmt.Errorf("invalid port number: %d > %d", port, maxPort)
	}
	return
}
