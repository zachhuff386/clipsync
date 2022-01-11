package utils

import (
	"fmt"
)

func LogError(err error) {
	fmt.Println(err.Error())
}
