package fwsmConfig

import (
	"fmt"
)

func warning(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	fmt.Println("")
}
