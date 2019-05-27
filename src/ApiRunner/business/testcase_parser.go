// testcase_parser.go
package business

import (
	"ApiRunner/models"
	"ApiRunner/utils"
	"log"

	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	toml "github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

func ParseTestCase(filePath string) (caseObj models.ICaseObj, err error) {
	if !utils.Exists(filePath) {
		err = errors.New(fmt.Sprintf(`testcase file not found,%s`, filePath))
		return
	}
	content, e := ioutil.ReadFile()
	if e != nil {
		err = e
		return
	}
	switch filepath.Ext(filePath) {
	case `.yaml`, `yml`:
		return parseYamlCase(content)
	case `.json`, `conf`:
		return parseJsonCase(content)
	case `.toml`, `.tml`:
		return parseTomlCase(content)
	}
}

func parseYamlCase(content []byte) (caseObj models.ICaseObj, err error) {
	m := make(map[string]interface{})
	err = yaml.Unmarshal([]byte(content), &m)
	if err != nil {
		log.Printf("parse yaml error: %v", err)
		return
	}
	fmt.Printf("--- yaml m:\n%v\n\n", m)
	caseObj, err = _parseTestCase(m)
	return
}

func parseJsonCase(content []byte) (caseObj models.ICaseObj, err error) {
	m := make(map[string]interface{})
	err = json.Unmarshal([]byte(content), &m)
	if err != nil {
		log.Printf("parse json error: %v", err)
		return
	}
	fmt.Printf("--- json m:\n%v\n\n", m)
	caseObj, err = _parseTestCase(m)
	return
}

func parseTomlCase(content []byte) (caseObj models.ICaseObj, err error) {
	m := make(map[string]interface{})
	err = toml.Unmarshal([]byte(content), &m)
	if err != nil {
		log.Printf("parse toml error: %v", err)
		return
	}
	fmt.Printf("--- toml m:\n%v\n\n", m)
	caseObj, err = _parseTestCase(m)
	return
}

func _parseTestCase(ts map[string]interface{}) (caseObj models.ICaseObj, err error) {
	return
}
