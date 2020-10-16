package conf

import "net/url"

var (
	config *CConfig
)

func init() {
	config = NewCConfig("config/config.json")
}

func GetConfigValues(index int) url.Values {
	return config.GetConfigValues(index)
}

func GetConfig(key string,index int) interface{} {
	return  config.GetConfig(key,index)
}

func GetConfigArray(key string) []interface{} {
	return config.GetConfigArray(key)
}