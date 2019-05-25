// testcase_parser_test.go
package business

import (
	"fmt"
	"testing"
)

func TestParseTestCase(t *testing.T) {
	t.Log(`TestParseTestCase`)
	o, e := ParseTestCase(`case.yaml`)
	fmt.Printf(`%v,%v`, o, e)

}
