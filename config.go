package longh

import (
	"encoding/json"
	"io/ioutil"
)

// LoadJSONConfig 读取配置文件 json格式
func LoadJSONConfig(filename string, v interface{}) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, v)
	if err != nil {
		return err
	}
	return nil
}
