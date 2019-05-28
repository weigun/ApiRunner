// testcase_parser_test.go
package business

import (
	"fmt"
	"testing"
)

func TestParseTestCase(t *testing.T) {
	t.Log(`TestParseTestCase`)
	o := ParseTestCase(`case.yaml`)
	fmt.Printf(`%v`, o)

}
