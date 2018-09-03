package utils

import (
	"bytes"
	"encoding/json"
	_ "fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"text/template"
	"time"
)

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

//func runLocalTasks(tasks []string) {
//	//单次执行用例
//	eng := NewEngine()
//	curFolder := getCwd()
//	var wg sync.WaitGroup
//	for i, v := range tasks {
//		casePath := filepath.Join(curFolder, "testcase", v+".json")
//		caseParser := NewCaseParser(casePath)
//		rn := NewRunner{caseParser.getCaseset()}
//		wg.Add(1)
//		go func(rn runner) {
//			defer wg.Done()  //TODO 可能需要放在safeRun下一行
//			rn.ready <- true //缓冲chan
//			eng.safeRun(rn)
//		}(rn) //copy rn
//	}
//	wg.Wait()
//	//	TODO 生成报告
//	log.Println("test done")
//}

//func getErr(api string, code int) string {
//	//请求错误码格式化
//	return fmt.Sprintf(`%s failed,StatusCode is %d`, api, code)
//}

type dataInterface interface {
	GetData() map[string]interface{}
}

func GetTemplate(_func *template.FuncMap) *template.Template {
	t := template.New("conf")
	if _func != nil {
		t.Funcs(*_func)
	}
	return t //不能放到全局或者通过闭包的方式，因为这个是携程不安全的
}

func Translate(tmpl *template.Template, tmplStr string, obj dataInterface) string {
	//将模板翻译
	wr := bytes.NewBufferString("")
	tmpl, err := tmpl.Parse(tmplStr)
	if err != nil {
		log.Fatalln(err)
	}
	if obj == nil {
		tmpl.Execute(wr, nil)
	} else {
		log.Println("===============", obj, obj.GetData())
		tmpl.Execute(wr, obj.GetData())
	}
	return wr.String()
}

//func translateValidata(data string, resp validation) string {
//	//将模板翻译
//	tmpl := getTemplate() //不能放到全局或者通过闭包的方式，因为这个是携程不安全的
//	wr := bytes.NewBufferString("")
//	tmpl, err := tmpl.Parse(data)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	switch resp.Body.(type) {
//	case respBodyMap:
//		tmpl.Execute(wr, resp.Body.(respBodyMap))
//	case respBodySlice:
//		tmpl.Execute(wr, resp.Body.(respBodySlice))
//	}
//	return wr.String()
//}

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
