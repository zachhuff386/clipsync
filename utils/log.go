package utils

import (
	"fmt"
)

func LogError(clientId string, err error) {
	fmt.Printf("[%s] %s\n", clientId, err.Error())
}
