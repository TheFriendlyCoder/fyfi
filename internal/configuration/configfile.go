package configuration

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type ConfigData struct {
	CacheFolder string
}

func NewConfigData(configFile string) (*ConfigData, error) {
	yfile, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	data := make(map[interface{}]interface{})
	err = yaml.Unmarshal(yfile, &data)
	if err != nil {
		return nil, err
	}

	var ok bool
	c := new(ConfigData)

	key := "cache_folder"
	c.CacheFolder, ok = data[key].(string)
	if !ok {
		return nil, fmt.Errorf("failed to parse %s from config file", key)
	}

	key = "port"
	var port int
	var temp interface{}
	temp, ok = data[key]
	if !ok {
		return nil, fmt.Errorf("unable to find key %s", key)
	}
	port, ok = temp.(int)
	fmt.Println(port)
	if !ok {
		return nil, fmt.Errorf("failed to parse %s from config file - incorrect data type %T", key, key)
	}
	return c, nil
}
