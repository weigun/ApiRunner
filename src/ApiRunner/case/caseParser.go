package testcase

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type translate interface {
	//转换接口
	conver(string) string        //将含有变量与表达式的模板翻译过来
	buildRequest() *http.Request //构造请求体
}

type casePaser struct {
	caseset *caseset
}

func NewCaseParser(casePth string) *casePaser {
	this := casePaser{}
	this.parse(casePth)
	return &this

}

func (this *casePaser) getCaseset() *caseset {
	return this.caseset
}

func (this *casePaser) parse(casePth string) {
	//{
	//    "name":"demo",
	//    "host":"http://10.104.225.242:9024",
	//    "headers":{
	//        "auth":"asdasda456789",
	//        "type":"json"
	//    },
	//    "globalVars":{
	//        "token":"{{getToken}}",
	//        "num":1
	//    },
	//}
	casePth = "D:\\test-area\\oschina\\ApiRunner\\src\\case\\demo.conf"
	fb, err := ioutil.ReadFile(casePth)
	if err != nil {
		fmt.Printf("%s\n", err)
		panic(err)
	}
	csMap := json2Map(fb)
	this.caseset.conf.name = csMap["name"].(string)
	this.caseset.conf.host = csMap["host"].(string)
	if _, ok := csMap["headers"].(map[string]string); ok {
		for k, v := range csMap["headers"].(map[string]string) {
			//增加公共头部
			this.caseset.conf.headers = append(this.caseset.conf.headers, header{k, v})
		}
	}
	if _, ok := csMap["globalVars"].(map[string]interface{}); ok {
		var h header
		for k, v := range csMap["globalVars"].(map[string]interface{}) {
			//增加用例集域的全局变量
			switch v.(type) {
			case int:
				h = variables{k, strconv.Itoa(v.(int))}
			case string:
				h = variables{k, v.(string)}
			case float64:
				h = variables{k, strconv.Itoa(v.(float64))}
			case []int:
				h = variables{k, v.([]int)}
			case []string:
				h = variables{k, v.([]string)}
			default:
				fmt.Printf("caseset globalVars not support type: %T", v)
				panic("caseset globalVars not support type")

			}
			this.caseset.conf.globalVars = append(this.caseset.conf.globalVars, h)
		}
	}
	//解析用例集配置完毕
	//开始解析用例
	for _, _case := range csMap["cases"].([]map[string]interface{}) {
		//    "cases":[
		//        {
		//            "name":"login",
		//            "api":"/api/user/userinfo",
		//            "method":"GET",
		//            "headers":{
		//                "cache":"true",
		//            },
		//            "params":{
		//                "time":1534991568.14434
		//            },
		//            "validate":[
		//                {
		//                    "op":"eq",
		//                    "source":"{{.body.code}}",
		//                    "verified":200
		//                },
		//                {
		//                    "op":"gt",
		//                    "source":"{{.body.data.num}}",
		//                    "verified":{{.num}}
		//                }
		//            ]
		//        }
		//    ]
		ci := caseItem{}
		ci.name = _case["name"].(string)
		ci.api = _case["api"].(string)
		ci.method = _case["method"].(string)
		if _, ok := _case["headers"].(map[string]string); ok {
			//先接入公共头部
			ci.headers = this.caseset.conf.headers
			for k, v := range _case["headers"].(map[string]string) {
				//增加用例头部,如果字段相同，会覆盖公共头部的字段
				index := ci.hasHeader(k)
				if index != -1 {
					ci.headers[index].val = v
				} else {
					ci.addHeader(header{k, v})
				}
			}
		}
		if _, ok := _case["params"].(map[string]interface{}); ok {
			//请求参数
			ci.params = _case["params"].(map[string]interface{})
		}
		for i, vali := range _case["validate"].([]map[string]string) {
			//验证条件
			c := condition{vali["op"], vali["source"], vali["verified"]}
			ci.addCondition(c)
		}
		this.caseset.addCaseItem(ci)
		fmt.Println("add caseItem:", ci)
	}
	//TODO 可能需要将string转为rune

}
