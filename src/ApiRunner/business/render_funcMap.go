// render_funcMap.go
package business

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Masterminds/sprig"
	// "github.com/zach-klippenstein/goregen"
)

//自定义一个funcMap，以sprig这个库为基础
var funcMap = sprig.TxtFuncMap()

//可以在这里编写自定义函数
func world() string {
	return `world`
}

func gen_email() string {
	/*email, err := regen.Generate("[a-z0-9]{4,32}")
	fmt.Println(`--------------------email:`, email)
	if err != nil {
		panic(err.Error())
	}
	return email
	*/
	myRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	randAlphaNum := funcMap[`randAlphaNum`].(func(count int) string)
	lenght := myRand.Intn(16)
	if lenght < 4 {
		lenght = 4
	}
	return fmt.Sprintf(`%s@%s.%s`,
		randAlphaNum(lenght),
		randAlphaNum(4),
		randAlphaNum(3))
}

//导入自定义函数
func init() {
	//导入
	funcMap[`world`] = world
	funcMap[`gen_email`] = gen_email
}
