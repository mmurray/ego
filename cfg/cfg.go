package cfg

import (
	"log"
	"path"
	"io/ioutil"
	"strings"
)

var cfgs = make(map[string]*ConfigMap)

func Parse(key string, filename string) {
	m := &ConfigMap{}
	m.Parse(filename)
	cfgs[key] = m
}

func ParseDir(dirname string) {
	dirlist, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Fatalf("Error reading %s: %s\n", dirname, err)
	}
	for _, f := range dirlist {
		filename := path.Join(dirname, f.Name())
		if f.IsDir() {
			ParseDir(filename)
		} else {
			Parse(strings.Replace(f.Name(), ".json", "", -1), filename)
		}
	}
}

func Get(key string) *ConfigMap {
	return cfgs[key]
}

