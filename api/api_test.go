package api_test

import "fmt"

func cvtErr(concrete, intf, method string) string {
	err := "interface conversion: %s is not %s: missing method %s"
	return fmt.Sprintf(err, concrete, intf, method)
}
