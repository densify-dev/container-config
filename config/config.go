package config

import (
	"fmt"
	rconf "github.com/densify-dev/retry-config/config"
	rconsts "github.com/densify-dev/retry-config/consts"
	"github.com/prometheus/common/model"
	"github.com/prometheus/common/sigv4"
	"strconv"
	"strings"
)

const (
	Empty = rconsts.Empty
)

type ClusterFilterParameters struct {
	Name        string         `yaml:"name"`
	Identifiers model.LabelSet `yaml:"identifiers,omitempty"`
}

type DensifyParameters struct {
	UrlConfig   *UrlConfig         `yaml:"url"`
	Endpoint    string             `yaml:"endpoint"`
	RetryConfig *rconf.RetryConfig `yaml:"retry,omitempty"`
}

type ProxyParameters struct {
	UrlConfig *UrlConfig `yaml:"url"`
	Auth      string     `yaml:"auth,omitempty"`
	Server    string     `yaml:"server,omitempty"`
	Domain    string     `yaml:"domain,omitempty"`
}

type ForwarderParameters struct {
	Densify *DensifyParameters `yaml:"densify"`
	Proxy   *ProxyParameters   `yaml:"proxy,omitempty"`
	Prefix  string             `yaml:"prefix,omitempty"`
}

type PrometheusParameters struct {
	UrlConfig   *UrlConfig         `yaml:"url"`
	BearerToken string             `yaml:"bearer_token,omitempty"`
	CaCertPath  string             `yaml:"ca_cert,omitempty"`
	SigV4Config *sigv4.SigV4Config `yaml:"sigv4,omitempty"`
	RetryConfig *rconf.RetryConfig `yaml:"retry,omitempty"`
}

type CollectionParameters struct {
	Include       map[string]bool `yaml:"include,omitempty"`
	Interval      string          `yaml:"interval"`
	IntervalSize  uint64          `yaml:"interval_size"`
	History       uint64          `yaml:"history"`
	HistoryInt    int             `yaml:"-"`
	Offset        uint64          `yaml:"offset"`
	OffsetInt     int             `yaml:"-"`
	SampleRate    uint64          `yaml:"sample_rate"`
	SampleRateSt  string          `yaml:"-"`
	NodeGroupList string          `yaml:"node_group_list"`
}

type Parameters struct {
	Forwarder  *ForwarderParameters       `yaml:"forwarder"`
	Prometheus *PrometheusParameters      `yaml:"prometheus"`
	Collection *CollectionParameters      `yaml:"collection"`
	Clusters   []*ClusterFilterParameters `yaml:"clusters"`
	Debug      bool                       `yaml:"debug"`
}

func merge(p *Parameters, pm *parameterMap) (newP *Parameters, err error) {
	if p == nil {
		pm.finalize()
		includes, _ := getIncludes(pm)
		cfp, _ := getClusterFilterParameters(pm)
		newP = &Parameters{
			Forwarder: &ForwarderParameters{
				Densify: &DensifyParameters{
					UrlConfig: getUrlConfig(pm, []string{densifyScheme, densifyHost, densifyPort, densifyUser, densifyPassword, densifyEncPassword}),
					Endpoint:  pm.stringValues[densifyEndpoint].v,
				},
				Proxy: &ProxyParameters{
					UrlConfig: getUrlConfig(pm, []string{proxyScheme, proxyHost, proxyPort, proxyUser, proxyPassword, proxyEncPassword}),
					Auth:      pm.stringValues[proxyAuth].v,
					Server:    pm.stringValues[proxyServer].v,
					Domain:    pm.stringValues[proxyDomain].v,
				},
				Prefix: pm.stringValues[filePrefix].v,
			},
			Prometheus: &PrometheusParameters{
				UrlConfig:   getUrlConfig(pm, []string{promScheme, promHost, promPort, promUser, promPassword, promPassword}),
				BearerToken: pm.stringValues[promToken].v,
				CaCertPath:  pm.stringValues[caCert].v,
			},
			Collection: &CollectionParameters{
				Include:       includes,
				Interval:      pm.stringValues[interval].v,
				IntervalSize:  pm.uint64Values[intervalSize].v,
				History:       pm.uint64Values[history].v,
				Offset:        pm.uint64Values[offset].v,
				SampleRate:    pm.uint64Values[sampleRate].v,
				NodeGroupList: pm.stringValues[nodeGroupList].v,
			},
			Clusters: []*ClusterFilterParameters{cfp},
			Debug:    pm.boolValues[debug].v,
		}
	} else {
		if pm.fromFile {
			err = fmt.Errorf("yaml configuration and properties configuration are mutually exclusive")
			return
		} else {
			newP = p
			newP.ensureStructs()
			// forwarder parameters:
			// 		densify parameters
			setValue(&newP.Forwarder.Densify.UrlConfig.Scheme, pm.stringValues, densifyScheme)
			setValue(&newP.Forwarder.Densify.UrlConfig.Host, pm.stringValues, densifyHost)
			setValue(&newP.Forwarder.Densify.UrlConfig.Port, pm.uint64Values, densifyPort)
			setValue(&newP.Forwarder.Densify.UrlConfig.Username, pm.stringValues, densifyUser)
			setValue(&newP.Forwarder.Densify.UrlConfig.Password, pm.stringValues, densifyPassword)
			setValue(&newP.Forwarder.Densify.UrlConfig.EncryptedPassword, pm.stringValues, densifyEncPassword)
			setValue(&newP.Forwarder.Densify.Endpoint, pm.stringValues, densifyEndpoint)
			// 		proxy parameters
			setValue(&newP.Forwarder.Proxy.UrlConfig.Scheme, pm.stringValues, proxyScheme)
			setValue(&newP.Forwarder.Proxy.UrlConfig.Host, pm.stringValues, proxyHost)
			setValue(&newP.Forwarder.Proxy.UrlConfig.Port, pm.uint64Values, proxyPort)
			setValue(&newP.Forwarder.Proxy.UrlConfig.Username, pm.stringValues, proxyUser)
			setValue(&newP.Forwarder.Proxy.UrlConfig.Password, pm.stringValues, proxyPassword)
			setValue(&newP.Forwarder.Proxy.UrlConfig.EncryptedPassword, pm.stringValues, proxyEncPassword)
			setValue(&newP.Forwarder.Proxy.Auth, pm.stringValues, proxyAuth)
			setValue(&newP.Forwarder.Proxy.Server, pm.stringValues, proxyServer)
			setValue(&newP.Forwarder.Proxy.Domain, pm.stringValues, proxyDomain)
			// 		prefix parameters
			setValue(&newP.Forwarder.Prefix, pm.stringValues, filePrefix)
			// prometheus parameters
			setValue(&newP.Prometheus.UrlConfig.Scheme, pm.stringValues, promScheme)
			setValue(&newP.Prometheus.UrlConfig.Host, pm.stringValues, promHost)
			setValue(&newP.Prometheus.UrlConfig.Port, pm.uint64Values, promPort)
			setValue(&newP.Prometheus.UrlConfig.Username, pm.stringValues, promUser)
			setValue(&newP.Prometheus.UrlConfig.Password, pm.stringValues, promPassword)
			setValue(&newP.Prometheus.BearerToken, pm.stringValues, promToken)
			setValue(&newP.Prometheus.CaCertPath, pm.stringValues, caCert)
			// collection parameters
			if includes, set := getIncludes(pm); set {
				newP.Collection.Include = includes
			}
			setValue(&newP.Collection.Interval, pm.stringValues, interval)
			setValue(&newP.Collection.IntervalSize, pm.uint64Values, intervalSize)
			setValue(&newP.Collection.History, pm.uint64Values, history)
			setValue(&newP.Collection.Offset, pm.uint64Values, offset)
			setValue(&newP.Collection.SampleRate, pm.uint64Values, sampleRate)
			setValue(&newP.Collection.NodeGroupList, pm.stringValues, nodeGroupList)
			// cluster parameter
			if cfp, set := getClusterFilterParameters(pm); set {
				newP.Clusters = append(newP.Clusters, cfp)
			}
			// debug parameter
			setValue(&newP.Debug, pm.boolValues, debug)
		}
	}
	err = newP.finalize()
	return
}

func (p *Parameters) ensureStructs() {
	if p.Forwarder == nil {
		p.Forwarder = &ForwarderParameters{}
	}
	if p.Forwarder.Densify == nil {
		p.Forwarder.Densify = &DensifyParameters{}
	}
	if p.Forwarder.Densify.UrlConfig == nil {
		p.Forwarder.Densify.UrlConfig = &UrlConfig{}
	}
	if p.Forwarder.Proxy == nil {
		p.Forwarder.Proxy = &ProxyParameters{}
	}
	if p.Forwarder.Proxy.UrlConfig == nil {
		p.Forwarder.Proxy.UrlConfig = &UrlConfig{}
	}
	if p.Prometheus == nil {
		p.Prometheus = &PrometheusParameters{}
	}
	if p.Prometheus.UrlConfig == nil {
		p.Prometheus.UrlConfig = &UrlConfig{}
	}
	if p.Collection == nil {
		p.Collection = &CollectionParameters{}
	}
}

// setValue can only be an issue if we have a parameter with default different from the type's zero value
// AND the zero value is a valid value - we do not have any parameter like this
func setValue[T comparable](target *T, vals values[T], name string) {
	if val, ok := vals[name]; ok {
		var zeroValue T
		if val.isSet || *target == zeroValue {
			*target = val.v
		}
	}
}

func getUrlConfig(pm *parameterMap, paramNames []string) *UrlConfig {
	return &UrlConfig{
		Scheme:            pm.stringValues[paramNames[0]].v,
		Host:              pm.stringValues[paramNames[1]].v,
		Port:              pm.uint64Values[paramNames[2]].v,
		Username:          pm.stringValues[paramNames[3]].v,
		Password:          pm.stringValues[paramNames[4]].v,
		EncryptedPassword: pm.stringValues[paramNames[5]].v,
	}
}

func getIncludes(pm *parameterMap) (m map[string]bool, set bool) {
	if val, ok := pm.stringValues[include]; ok {
		set = val.isSet
		vals := strings.Split(strings.ToLower(val.v), Comma)
		m = make(map[string]bool, len(vals))
		for _, v := range vals {
			m[v] = true
		}
	}
	return
}

func getClusterFilterParameters(pm *parameterMap) (cfp *ClusterFilterParameters, set bool) {
	if val, ok := pm.stringValues[clusterName]; ok {
		set = val.isSet
		cfp = &ClusterFilterParameters{Name: val.v}
	}
	return
}

func (p *Parameters) finalize() (err error) {
	p.Collection.HistoryInt = int(p.Collection.History)
	p.Collection.OffsetInt = int(p.Collection.Offset)
	p.Collection.SampleRateSt = strconv.FormatUint(p.Collection.SampleRate, 10)
	if err = p.Forwarder.Densify.UrlConfig.finalize(); err != nil {
		return
	}
	if err = p.Forwarder.Densify.RetryConfig.Validate(); err != nil {
		return
	}
	if err = p.Forwarder.Proxy.UrlConfig.finalize(); err != nil {
		return
	}
	if err = p.Prometheus.UrlConfig.finalize(); err != nil {
		return
	}
	err = p.Prometheus.RetryConfig.Validate()
	return
}
