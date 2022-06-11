package configuration

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type ConfigData struct {
	// Path to the folder where cached packages will be stored
	// This will also be the location of the metadata database
	// containing reference information about the cached data
	CacheFolder string
}

// NewConfigData constructs an instance of the ConfigData
// struct, populating the contents from data loaded from
// a yaml-formatted configuration file. The file format is
// expected to look as follows:
//
//		cache_folder: <absolute_path_to_cache_folder>
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

	return c, nil
}
