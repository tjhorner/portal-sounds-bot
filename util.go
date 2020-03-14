package main

import (
	"fmt"
)

func Ellipsize(str string, max int) string {
	if len(str) <= max {
		return str
	}

	return fmt.Sprintf("%s…", str[0:max-2])
}
