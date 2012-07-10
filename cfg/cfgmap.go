package cfg

import (
	"encoding/json"
	"io/ioutil"
)

type ConfigMap map[string]interface{}

func (m *ConfigMap) Parse(filename string) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, m)
	if err != nil {
		panic(err)
	}
}