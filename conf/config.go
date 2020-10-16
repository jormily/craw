package conf

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"
)

type CConfig struct {
	data map[string]interface{}
}

func NewCConfig(filename string) *CConfig {
	this := new(CConfig)
	this.data = make(map[string]interface{})
	err := this.Init(filename)
	if err != nil {
		log.Fatal(err.Error())
	}

	return this
}

func (this *CConfig)Init(filename string) error {
	cnf, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(cnf, &this.data)
	if err != nil {
		return err
	}

	return nil
}

func (this *CConfig)GetConfigString(key string) string {
	if config, ok := this.data[key]; ok {
		if configv, ok := config.(string); ok {
			return configv
		}
	}
	return ""
}

func (this *CConfig)GetConfigArray(key string) []interface{} {
	if config, ok := this.data[key]; ok {
		if configv, ok := config.([]interface{}); ok {
			return configv
		}
	}
	return nil
}

func (this *CConfig)GetConfig(key string,index int) interface{} {
	index = index - 1
	if config, ok := this.data[key]; ok {
		if configv, ok := config.([]interface{});ok {
			if configv[index] != nil {
				return configv[index]
			}
		}
	}
	return nil
}

func (this *CConfig)GetConfigValues(index int) url.Values {
	var cnf = this.GetConfig("values", index)
	if cnf == nil {
		return nil
	}

	if values,ok := cnf.(map[string]interface{}); ok {
		val := url.Values{}
		for k,v := range values {
			val.Set(k ,v.(string))
		}

		return val
	}

	return nil
}
