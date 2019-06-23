package business

import (
	"log"

	assertLib "github.com/smartystreets/assertions"
)

const (
	SUCCESS = ``
)

type assertion func(actual interface{}, expected ...interface{}) bool

func Equal(actual interface{}, expected ...interface{}) bool {
	ret := assertLib.ShouldEqual(actual, expected...)
	log.Println(`Assert Equal:`, ret)
	return ret == SUCCESS
}

func GreaterThan(actual interface{}, expected ...interface{}) bool {
	ret := assertLib.ShouldBeGreaterThan(actual, expected...)
	log.Println(`Assert GreaterThan:`, ret)
	return ret == SUCCESS
}

func LessThan(actual interface{}, expected ...interface{}) bool {
	ret := assertLib.ShouldBeLessThan(actual, expected...)
	log.Println(`Assert ShouldBeLessThan:`, ret)
	return ret == SUCCESS
}

func LessThanOrEqual(actual interface{}, expected ...interface{}) bool {
	ret := assertLib.ShouldBeLessThanOrEqualTo(actual, expected...)
	log.Println(`Assert ShouldBeLessThanOrEqualTo:`, ret)
	return ret == SUCCESS
}

func GreaterThanOrEqual(actual interface{}, expected ...interface{}) bool {
	ret := assertLib.ShouldBeGreaterThanOrEqualTo(actual, expected...)
	log.Println(`Assert ShouldBeGreaterThanOrEqualTo:`, ret)
	return ret == SUCCESS
}

func Contain(actual interface{}, expected ...interface{}) bool {
	ret := assertLib.ShouldContain(actual, expected...)
	log.Println(`Assert ShouldContain:`, ret)
	return ret == SUCCESS
}

func HaveLength(actual interface{}, expected ...interface{}) bool {
	ret := assertLib.ShouldHaveLength(actual, expected...)
	log.Println(`Assert ShouldHaveLength:`, ret)
	return ret == SUCCESS
}

func So(actual interface{}, assert assertion, expected ...interface{}) bool {
	isok, result := assertLib.So(actual, assert, expected...)
	if !isok {
		log.Println(`So failed.result is `, result)
	}
	return isok
}

var assertMap = map[string]assert{
	`eq`:    Equal,
	`equal`: Equal,
	`gt`:    GreaterThan,
	`lt`:    LessThan,
	`le`:    LessThanOrEqual,
	`ge`:    GreaterThanOrEqual,
	`in`:    Contain,
	`len`:   HaveLength,
}

func getAssertByOp(op string) assertion {
	if ok, assertFun := assertMap[op]; ok {
		return assertFun
	}
	log.Panicln(`unknow assertion:`, op)
}
