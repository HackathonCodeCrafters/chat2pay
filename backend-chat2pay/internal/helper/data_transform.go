package helper

import "fmt"

func InterfaceToString(v interface{}) string {
	return fmt.Sprintf(`%v`, v)
}
