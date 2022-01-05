package util

import (
	"math/rand"
	"strings"
	"time"
)

func RandomString(size int) string {
	str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

	result := make([]string, size)
	rand.Seed(time.Now().UnixNano())

	for i := range result {
		result[i] = string(str[rand.Intn(len(str))])
	}

	return strings.Join(result, "")
}
