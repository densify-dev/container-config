package config

import (
	"fmt"
	nnet "github.com/densify-dev/net-utils/network"
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
	Slash                   = "/"
	Http                    = "http"
	Https                   = Http + "s"
	DefaultHttpPort  uint64 = 80
	DefaultHttpsPort uint64 = 443
	IgnorePort       uint64 = 99999 // 0 is a valid port, need another invalid value indicating "ignore me"
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
	hostElems := strings.SplitN(uc.Host, Slash, 2)
	var h string
	if omitPort(sc, uc.Port) {
		h = hostElems[0]
	} else {
		var p nnet.Port
		if p, err = nnet.NewPort(uc.Port); err == nil {
			h = p.Addr(hostElems[0])
		} else {
			return
		}
	}
	u := &url.URL{Scheme: sc, Host: h}
	if len(hostElems) > 1 {
		u.Path = Slash + hostElems[1]
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

func omitPort(scheme string, port uint64) bool {
	return port == IgnorePort ||
		(scheme == Http) && (port == DefaultHttpPort) ||
		(scheme == Https) && (port == DefaultHttpsPort)
}
