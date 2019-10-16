//hello_test.go
package hello

import "testing"

func Testhellogolang(t *testing.T)  {
	if hellogolang() == "Hello, World!" {
		t.Log("Test Passed")
	} else {
		t.Error("Test Failed")
	}
}