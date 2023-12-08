package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

func ReadConfig() (p *Parameters, err error) {
	pm := initParameterMap()
	var fc *fileConfig
	if fc, err = pm.populate(); err != nil {
		return
	}
	if fc.configType() == hierarchyType {
		if p, err = readParams(fc.path()); err != nil {
			return
		}
	}
	p, err = merge(p, pm)
	return
}

type configFileType int

const (
	unknownType configFileType = iota
	mapType
	hierarchyType
)

var fileTypeMapping = map[string]configFileType{
	"yaml":       hierarchyType,
	"yml":        hierarchyType,
	"properties": mapType,
	"props":      mapType,
}

type fileConfig struct {
	dir  string
	file string
	typ  string
}

func (fc *fileConfig) configType() (cft configFileType) {
	var f bool
	if fc != nil {
		if cft, f = fileTypeMapping[fc.typ]; !f {
			ext := filepath.Ext(fc.file)
			cft, f = fileTypeMapping[ext[1:]]
		}
	}
	if !f {
		cft = unknownType
	}
	return
}

func (fc *fileConfig) path() string {
	file := fc.file
	if ext := filepath.Ext(file); ext == Empty && fc.typ != Empty {
		file = strings.Join([]string{file, fc.typ}, Dot)
	}
	return filepath.Join(fc.dir, file)
}

func readParams(config string) (p *Parameters, err error) {
	var data []byte
	if data, err = os.ReadFile(config); err != nil {
		return
	}
	p = &Parameters{}
	if err = yaml.Unmarshal(data, p); err != nil {
		p = nil
	}
	return
}
