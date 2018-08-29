package engine

//import (
//	tils "ApiRunner/utils"
//	"fmt"
//	"log"
//	"strconv"
//)

//type respBodyMap struct {
//	//返回是map类型的
//	Code int
//	Msg  string
//	Data map[string]interface{}
//}

//type respBodySlice struct {
//	//返回是列表类型的
//	Code int
//	Msg  string
//	Data []map[string]interface{}
//}

//type validation struct {
//	//验证体
//	//定义验证数据的结构，可通过模板翻译，用来做结果对比
//	Code int
//	Body interface{}
//}

//func validate(resp Response, conds []condition) {
//	//TODO 结果需要给报告
//	_ := resp
//	var ntr interface{}
//	contentMap := utils.Json2Map([]byte(resp.Content))
//	switch contentMap["data"].(type) {
//	case map[string]interface{}:
//		ntr = respBodyMap{resp.Code, contentMap["msg"].(string), contentMap}
//	case []map[string]interface{}:
//		ntr = respBodySlice{resp.Code, contentMap["msg"].(string), contentMap}
//	default:
//		log.Fatal("contentMap not support type")
//	}
//	for i, cond := range conds {
//		fmt.Println(cond)
//		ret := translateValidata(cond.source, validation{resp.Code, ntr})
//		if cond.operation == eq {
//			switch cond.verified.(type) {
//			case string:
//				if cond.verified == ret {
//					fmt.Println("cond.verified is string")
//					fmt.Println("compare:", cond.verified, ret, verified == ret)
//				}
//			case int:
//				if cond.verified == strconv.Atoi(ret) {
//					fmt.Println("cond.verified is int")
//					fmt.Println("compare:", cond.verified, ret, verified == ret)
//				}
//			default:
//				fmt.Println("cond.verified is not string or int")
//			}
//		}
//	}

//}
