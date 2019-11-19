package utils

import (
	// "bytes"
	// "encoding/json"
	// "io/ioutil"
	"crypto/md5"
	"encoding/hex"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Map2Json(m map[string]interface{}) string {
	jsonStr, err := json.Marshal(m)
	if err != nil {
		return "{}"
	}
	return string(jsonStr)
}

func Json2Map(js []byte) map[string]interface{} {
	var mapResult = make(map[string]interface{})
	//使用 json.Unmarshal(data []byte, v interface{})进行转换,返回 error 信息
	if err := json.Unmarshal(js, &mapResult); err != nil {
		return mapResult
	}
	return mapResult
}

func GetCwd() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func ToNumber(a interface{}) interface{} {
	switch a.(type) {
	case int, float64:
		return a
	case string:
		a := a.(string)
		i, err := strconv.ParseInt(a, 10, 64)
		if err != nil {
			i, err := strconv.ParseFloat(a, 64)
			if err != nil {
				return nil
			}
			return i
		}
		return i
	default:
		return nil
	}
}

func GetDateTime() string {
	timeStamp := time.Now().Unix()
	return time.Unix(timeStamp, 0).Format("20060102_150405")
}

func Now4ms() int64 {
	return time.Now().UnixNano() / 1e6
}

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
