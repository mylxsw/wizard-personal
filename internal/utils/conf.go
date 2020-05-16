package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

// LoadJSONFile 加载 Json 配置文件
func LoadJSONFile(filename string, res interface{}) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.Wrapf(err, "load file %s failed", filename)
	}

	if err := json.Unmarshal(content, res); err != nil {
		return errors.Wrapf(err, "parse file %s failed", filename)
	}

	return nil
}

// WriteJSONFile 将json写入文件
func WriteJSONFile(filename string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "can not serialize data to json")
	}

	var output bytes.Buffer
	if err := json.Indent(&output, jsonData, "", "    "); err != nil {
		return err
	}

	return ioutil.WriteFile(filename, output.Bytes(), os.ModePerm)
}

// FileExist 判断文件是否存在
func FileExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
