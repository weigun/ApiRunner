package validation

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"text/template"
)

func toNumber(a interface{}) interface{} {
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

type resp struct {
	Body map[string]interface{}
}

func getTemplate() *template.Template {
	t := template.New("conf")
	t.Funcs(template.FuncMap{"show": show})
	return t //不能放到全局或者通过闭包的方式，因为这个是携程不安全的
}

func show() string {
	return "success!!!!!!!!!!!!!"
}
func translate(data string, r resp) string {
	//将模板翻译
	log.Println(data, r)
	tmpl := getTemplate()
	wr := bytes.NewBufferString("")
	tmpl, err := tmpl.Parse(data)
	if err != nil {
		log.Fatalln(err)
	}
	tmpl.Execute(wr, r)
	return wr.String()
}

func _main() {
	fmt.Println("Hello World")
	for _, a := range []interface{}{"23.12", 123, 0.25, -1, "-100", "asd", nil} {
		b := toNumber(a)
		fmt.Println(a, b)
	}
	m := make(map[string]interface{})
	m["code"] = 200
	m["msg"] = "ok"
	m["data"] = make(map[string]interface{})
	m["data"].(map[string]interface{})["token name"] = "#$%^&*(asdadasdasdasd52454)"
	m["data"].(map[string]interface{})["firstTime"] = 1
	var l []map[string]interface{}
	l = append(l, map[string]interface{}{"name": "weigun"})
	l = append(l, map[string]interface{}{"name": "wgg"})
	l = append(l, map[string]interface{}{"name": "gun"})
	m["data"].(map[string]interface{})["list"] = l
	r := resp{m}
	ret := translate(`{{show}} resp is {{.Body.code}},{{.Body.msg}},{{index .Body.data "token name"}},{{.Body.data.firstTime}},{{index .Body.data.list 1 "name"}}`, r)
	log.Println(ret)
}
