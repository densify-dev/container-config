package config

import (
	"fmt"
	"github.com/go-viper/encoding/javaproperties"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"strings"
)

// keys
const (
	debug              = "debug"
	configDir          = "config_dir"
	configFile         = "config_file"
	configType         = "config_type"
	clusterName        = "cluster_name"
	promScheme         = "prometheus_protocol"
	promHost           = "prometheus_address"
	promPort           = "prometheus_port"
	promUser           = "prometheus_user"
	promPassword       = "prometheus_password"
	promToken          = "prometheus_oauth_token"
	caCert             = "ca_certificate"
	include            = "include_list"
	nodeGroupList      = "node_group_list"
	interval           = "interval"
	intervalSize       = "interval_size"
	sampleRate         = "sample_rate"
	history            = "history"
	offset             = "offset"
	densifyScheme      = "protocol"
	densifyHost        = "host"
	densifyPort        = "port"
	densifyEndpoint    = "endpoint"
	densifyUser        = "user"
	densifyPassword    = "password"
	densifyEncPassword = "epassword"
	proxyScheme        = "proxyprotocol"
	proxyHost          = "proxyhost"
	proxyPort          = "proxyport"
	proxyAuth          = "proxyauth"
	proxyServer        = "proxyserver"
	proxyDomain        = "proxydomain"
	proxyUser          = "proxyuser"
	proxyPassword      = "proxypassword"
	proxyEncPassword   = "eproxypassword"
	filePrefix         = "prefix"
)

// default values as consts
const (
	defConfigDir              = "./config"
	defConfigFile             = "config"
	defConfigType             = "properties"
	defPromScheme             = Http
	defPromPort        uint64 = 9090
	defInclude                = "container,node,cluster,nodegroup,quota"
	defNodeGroupList          = "label_karpenter_sh_nodepool,label_cloud_google_com_gke_nodepool,label_eks_amazonaws_com_nodegroup,label_agentpool,label_pool_name,label_alpha_eksctl_io_nodegroup_name,label_kops_k8s_io_instancegroup"
	defInterval               = "hours"
	defIntervalSize    uint64 = 1
	defSampleRate      uint64 = 5
	defHistory         uint64 = 1
	defDensifyPort            = DefaultHttpsPort
	defDensifyScheme          = Https
	defDensifyHost            = "localhost"
	defDensifyEndpoint        = "/api/v2/"
	defProxyPort              = DefaultHttpsPort
	defProxyAuth              = "Basic"
)

// default values as vars using Go's default values
var (
	defOffset uint64
	defDebug  bool
)

type valueSpec struct {
	name      string
	shorthand string
	usage     string
	v         *viper.Viper
}

type pflagFunc[T comparable] func(*T, string, string, T, string)
type getFunc[T comparable] func(*viper.Viper, string) T

type value[T comparable] struct {
	spec  *valueSpec
	v     T
	defV  T
	pf    pflagFunc[T]
	gf    getFunc[T]
	isSet bool
}

type values[T comparable] map[string]*value[T]

type parameterMap struct {
	keys         map[string]bool
	stringValues values[string]
	uint64Values values[uint64]
	boolValues   values[bool]
	fromFile     bool
}

func initParameterMap() *parameterMap {
	pm := &parameterMap{
		keys:         make(map[string]bool),
		stringValues: make(values[string]),
		uint64Values: make(values[uint64]),
		boolValues:   make(values[bool]),
	}
	// config file parameters
	_ = pm.addStringValue(configDir, "l", "config file parent directory", Empty, defConfigDir)
	_ = pm.addStringValue(configFile, "f", "config file name (without extension)", Empty, defConfigFile)
	_ = pm.addStringValue(configType, "y", "config file type", Empty, defConfigType)
	// debug parameter
	_ = pm.addBoolValue(debug, "d", "enable debug-level logging", Empty, defDebug)
	// single cluster parameter
	_ = pm.addStringValue(clusterName, "c", "cluster name", Empty, Empty)
	// prometheus parameters
	_ = pm.addStringValue(promScheme, "s", "prometheus scheme", Empty, defPromScheme)
	_ = pm.addStringValue(promHost, "a", "prometheus host", Empty, Empty)
	_ = pm.addUint64Value(promPort, "p", "prometheus port", Empty, defPromPort)
	_ = pm.addStringValue(promUser, "u", "prometheus basic auth user - value or filename", Empty, Empty)
	_ = pm.addStringValue(promPassword, "w", "prometheus basic auth password - value or filename", Empty, Empty)
	_ = pm.addStringValue(promToken, "t", "prometheus oauth token - value or filename", Empty, Empty)
	_ = pm.addStringValue(caCert, "x", "path to CA certificate (may be required to pass certificate validation)", Empty, Empty)
	// collection parameters
	_ = pm.addStringValue(include, "n", "comma-separated list of data to include in collection: cluster, node, container, nodegroup, quota", Empty, defInclude)
	_ = pm.addStringValue(nodeGroupList, "g", "comma-separated list of label names to check for building node groups", Empty, defNodeGroupList)
	_ = pm.addStringValue(interval, "k", "interval unit - days/hours/minutes", Empty, defInterval)
	_ = pm.addUint64Value(intervalSize, "i", "interval size to be used for querying - last interval size of interval unit of data", Empty, defIntervalSize)
	_ = pm.addUint64Value(sampleRate, "r", "rate of sample points to collect (1 sample every sample rate in minutes)", Empty, defSampleRate)
	_ = pm.addUint64Value(history, "h", "time to go back for data collection, works with the interval and interval size settings", Empty, defHistory)
	_ = pm.addUint64Value(offset, "o", "amount of units (based on interval value) to offset the data collection backwards in time", Empty, defOffset)
	// forwarder parameters
	// 		Densify parameters
	_ = pm.addStringValue(densifyScheme, "S", "densify scheme", forwarderEnvPrefix, defDensifyScheme)
	_ = pm.addStringValue(densifyHost, "H", "densify host", forwarderEnvPrefix, defDensifyHost)
	_ = pm.addUint64Value(densifyPort, "P", "densify port", forwarderEnvPrefix, defDensifyPort)
	_ = pm.addStringValue(densifyEndpoint, "N", "densify endpoint", forwarderEnvPrefix, defDensifyEndpoint)
	_ = pm.addStringValue(densifyUser, "U", "densify user - value or filename", forwarderEnvPrefix, Empty)
	_ = pm.addStringValue(densifyPassword, "W", "densify password - value or filename", forwarderEnvPrefix, Empty)
	_ = pm.addStringValue(densifyEncPassword, "E", "encrypted densify password - value or filename", forwarderEnvPrefix, Empty)
	// 		proxy parameters
	_ = pm.addStringValue(proxyScheme, "T", "proxy scheme", forwarderEnvPrefix, Empty)
	_ = pm.addStringValue(proxyHost, "G", "proxy host", forwarderEnvPrefix, Empty)
	_ = pm.addUint64Value(proxyPort, "Q", "proxy port", forwarderEnvPrefix, defProxyPort)
	_ = pm.addStringValue(proxyAuth, "A", "proxy auth", forwarderEnvPrefix, defProxyAuth)
	_ = pm.addStringValue(proxyServer, "R", "proxy server", forwarderEnvPrefix, Empty)
	_ = pm.addStringValue(proxyDomain, "D", "proxy domain", forwarderEnvPrefix, Empty)
	_ = pm.addStringValue(proxyUser, "V", "proxy user - value or filename", forwarderEnvPrefix, Empty)
	_ = pm.addStringValue(proxyPassword, "X", "proxy password - value or filename", forwarderEnvPrefix, Empty)
	_ = pm.addStringValue(proxyEncPassword, "F", "encrypted proxy password - value or filename", forwarderEnvPrefix, Empty)
	// 		forwarder parameters
	_ = pm.addStringValue(filePrefix, "I", "zip file prefix", forwarderEnvPrefix, Empty)

	return pm
}

func getString(v *viper.Viper, key string) string {
	return v.GetString(key)
}

func (pm *parameterMap) addStringValue(name, shorthand, usage string, envPrefix string, defV string) error {
	return addValue(pm.keys, pm.stringValues, name, shorthand, usage, envPrefix, defV, pflag.StringVarP, getString)
}

func getUint64(v *viper.Viper, key string) uint64 {
	return v.GetUint64(key)
}

func (pm *parameterMap) addUint64Value(name, shorthand, usage string, envPrefix string, defV uint64) error {
	return addValue(pm.keys, pm.uint64Values, name, shorthand, usage, envPrefix, defV, pflag.Uint64VarP, getUint64)
}

func getBool(v *viper.Viper, key string) bool {
	return v.GetBool(key)
}

func (pm *parameterMap) addBoolValue(name, shorthand, usage string, envPrefix string, defV bool) error {
	return addValue(pm.keys, pm.boolValues, name, shorthand, usage, envPrefix, defV, pflag.BoolVarP, getBool)
}

func addValue[T comparable](keys map[string]bool, vals values[T], name, shorthand, usage string, envPrefix string, defV T, pf pflagFunc[T], gf getFunc[T]) error {
	if keys[name] {
		return fmt.Errorf("duplicate key %s", name)
	}
	v, ok := vipersByPrefix[envPrefix]
	if !ok {
		return fmt.Errorf("invalid env prefix %s", envPrefix)
	}
	val := &value[T]{
		spec: &valueSpec{
			name:      name,
			shorthand: shorthand,
			usage:     usage,
			v:         v,
		},
		defV: defV,
		pf:   pf,
		gf:   gf,
	}
	vals[name] = val
	keys[name] = true
	return nil
}

const (
	forwarderEnvPrefix = "densify"
)

var envPrefixes = []string{Empty, forwarderEnvPrefix}
var vipers, vipersByPrefix = initVipers()

func initVipers() (vs []*viper.Viper, m map[string]*viper.Viper) {
	l := len(envPrefixes)
	vs = make([]*viper.Viper, l)
	m = make(map[string]*viper.Viper, l)
	codecRegistry := viper.NewCodecRegistry()
	codecRegistry.RegisterCodec(defConfigType, &javaproperties.Codec{})
	for i, envPrefix := range envPrefixes {
		v := viper.NewWithOptions(viper.WithCodecRegistry(codecRegistry))
		v.SetTypeByDefaultValue(true)
		v.SetEnvPrefix(envPrefix)
		v.AutomaticEnv()
		vs[i] = v
		m[envPrefix] = v
	}
	return
}

func (pm *parameterMap) populate() (fc *fileConfig, err error) {
	if err = populateValues(pm.stringValues); err == nil {
		if err = populateValues(pm.uint64Values); err == nil {
			err = populateValues(pm.boolValues)
		}
	}
	if err != nil {
		return
	}
	pflag.Parse()
	for _, v := range vipers {
		if err = v.BindPFlags(pflag.CommandLine); err != nil {
			return
		}
	}
	// the meta config (config path, filename and type) are only available at the first Viper instance
	fc = getFileConfig(vipers[0])
	for _, v := range vipers {
		if fc.configType() == mapType {
			v.SetConfigName(fc.file)
			v.SetConfigType(fc.typ)
			v.AddConfigPath(fc.dir)
			e := v.ReadInConfig()
			pm.fromFile = e == nil
		}
	}
	// resolve the values
	resolve(pm.stringValues)
	resolve(pm.uint64Values)
	resolve(pm.boolValues)
	return
}

func populateValues[T comparable](vals values[T]) error {
	// default values are used as defaults for pflag - if we'd call v.SetDefault(), then IsSet() for
	// that key will always return true
	for key, val := range vals {
		val.pf(&val.v, val.spec.name, val.spec.shorthand, val.defV, val.spec.usage)
		if err := val.spec.v.BindEnv([]string{key}...); err != nil {
			return err
		}
	}
	return nil
}

func getFileConfig(v *viper.Viper) *fileConfig {
	return &fileConfig{
		dir:  v.GetString(configDir),
		file: v.GetString(configFile),
		typ:  strings.ToLower(v.GetString(configType)),
	}
}

func resolve[T comparable](vals values[T]) {
	for key, val := range vals {
		val.v = val.gf(val.spec.v, key)
		val.isSet = val.spec.v.IsSet(key)
	}
}

func (pm *parameterMap) finalize() {
	if val, f := pm.stringValues[clusterName]; !f || val == nil || val.v == Empty {
		pm.stringValues[clusterName] = pm.stringValues[promHost]
	}
}
