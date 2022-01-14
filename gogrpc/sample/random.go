package sample

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomBrand(brand ...string) string {
	length := len(brand)
	b := brand[rand.Intn(length)]
	return b
}

func RandomName(name ...string) string {
	length := len(name)
	n := name[rand.Intn(length)]
	return n
}