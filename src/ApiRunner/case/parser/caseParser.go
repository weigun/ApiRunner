package parse

import (
	caseMod "ApiRunner/case"
	varsMgr "ApiRunner/manager/bucket/caseVariables"
	utils "ApiRunner/utils"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"log"
	_ "net/http"
	"strconv"
	"strings"
)

type casePaser struct {
	Caseset *caseMod.Caseset
	Uid     uint32
}

func NewCaseParser(casePth string) *casePaser {
	this := casePaser{Caseset: caseMod.NewCaseset()}
	this.parse(casePth)
	this.Uid = crc32.ChecksumIEEE([]byte(this.Caseset.Conf.Name))
	for _, v := range this.Caseset.Conf.GlobalVars {
		varsMgr.SetVar(this.Uid, v.Name, v.Conver())
		log.Println("add GlobalVars:", this.Uid, v.Name, v.Conver())
	}
	return &this

}

func (this *casePaser) GetCaseset() *caseMod.Caseset {
	return this.Caseset
}

func (this *casePaser) GetUid() uint32 {
	return this.Uid
}

func (this *casePaser) parse(casePath string) {
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
	//	casePath = `D:\test-area\github\ApiRunner\src\ApiRunner\case\demo.conf`
	log.Println("casePath:", casePath)
	fb, err := ioutil.ReadFile(casePath)
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	csMap := utils.Json2Map(fb)
	//	fmt.Println(csMap)
	this.Caseset.Conf.Name = csMap["name"].(string)
	this.Caseset.Conf.Host = csMap["host"].(string)
	if _, ok := csMap["headers"]; ok {
		for k, v := range csMap["headers"].(map[string]interface{}) {
			//增加公共头部
			this.Caseset.Conf.Headers = append(this.Caseset.Conf.Headers, caseMod.Header{k, v.(string)})
		}
		//		fmt.Println("hehe++++++++++++++++", this.Caseset.Conf.Headers)
	}
	if _, ok := csMap["globalVars"].(map[string]interface{}); ok {
		var h caseMod.Variables
		for k, v := range csMap["globalVars"].(map[string]interface{}) {
			//增加用例集域的全局变量
			switch v.(type) {
			case int:
				h = caseMod.Variables{k, strconv.Itoa(v.(int))}
			case string:
				h = caseMod.Variables{k, v.(string)}
			case float64:
				h = caseMod.Variables{k, strconv.FormatFloat(v.(float64), 'f', 3, 64)}
			case []int:
				h = caseMod.Variables{k, v.([]int)}
			case []string:
				h = caseMod.Variables{k, v.([]string)}
			default:
				log.Printf("caseset globalVars not support type: %T", v)
				panic("caseset globalVars not support type")

			}
			this.Caseset.Conf.GlobalVars = append(this.Caseset.Conf.GlobalVars, h)
		}
	}
	//解析用例集配置完毕
	//开始解析用例
	//	fmt.Printf(">>>>>>>>>>>%T\n", csMap["cases"])
	for _, _case := range csMap["cases"].([]interface{}) {
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
		_case := _case.(map[string]interface{})
		ci := caseMod.CaseItem{}
		ci.Name = _case["name"].(string)
		ci.Api = strings.Join([]string{this.Caseset.Conf.Host, _case["api"].(string)}, "")
		ci.Method = _case["method"].(string)
		if _, ok := _case["headers"]; ok {
			//先接入公共头部
			ci.Headers = this.Caseset.Conf.Headers
			for k, v := range _case["headers"].(map[string]interface{}) {
				//增加用例头部,如果字段相同，会覆盖公共头部的字段
				index := ci.HasHeader(k)
				if index != -1 {
					ci.Headers[index].Val = v.(string)
				} else {
					ci.AddHeader(caseMod.Header{k, v.(string)})
				}
			}
		} else {
			//直接用公共header
			//			fmt.Println("enter---------------------")
			ci.Headers = this.Caseset.Conf.Headers
			//			fmt.Println(ci.Headers, this.Caseset.Conf.Headers)
		}
		if _, ok := _case["params"].(map[string]interface{}); ok {
			//请求参数
			ci.Params.Params = _case["params"].(map[string]interface{})
		}
		//导出变量
		if _, ok := _case["export"].(map[string]interface{}); ok {
			var ev caseMod.Variables
			for k, v := range _case["export"].(map[string]interface{}) {
				//增加用例变量
				switch v.(type) {
				case int:
					ev = caseMod.Variables{k, strconv.Itoa(v.(int))}
				case string:
					ev = caseMod.Variables{k, v.(string)}
				case float64:
					ev = caseMod.Variables{k, strconv.FormatFloat(v.(float64), 'f', 3, 64)}
				case []int:
					ev = caseMod.Variables{k, v.([]int)}
				case []string:
					ev = caseMod.Variables{k, v.([]string)}
				default:
					log.Printf("caseset exportVars not support type: %T", v)
					panic("caseset exportVars not support type")

				}
				//				this.Caseset.Conf.GlobalVars = append(this.Caseset.Conf.GlobalVars, h)
				ci.AddExportVar(ev)
			}
		}
		for _, vali := range _case["validate"].([]interface{}) {
			//验证条件
			vali := vali.(map[string]interface{})
			c := caseMod.Condition{vali["op"].(string), vali["source"].(string), vali["verified"].(string)}
			ci.AddCondition(c)
		}
		this.Caseset.AddCaseItem(ci)
		fmt.Println("add caseItem:", ci)
	}
	//TODO 可能需要将string转为rune
	//TODO 用例需要池化，全局都可以能拿得到
	log.Println("===============", this.Caseset)

}

func pretreatment(casePath string) []byte {
	// 用例数据预处理
	// 先将全局变量实例化，这样方便后来的变量引用
	log.Println("casePath:", casePath)
	fb, err := ioutil.ReadFile(casePath)
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	log.Println("case data pretreating....")
	_ = fb
	return []byte{}
}
