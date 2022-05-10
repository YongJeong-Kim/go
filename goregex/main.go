package main

import (
	"fmt"
	"regexp"
)

func NumberOnly(input string) (bool, error) {
	pattern := "^[0-9]+$"
	ok, err := regexp.MatchString(pattern, input)
	if err != nil {
		return false, nil
	}
	return ok, nil
}

func LowerCaseOnly(input string) (bool, error) {
	pattern := "^[a-z]+$"
	ok, err := regexp.MatchString(pattern, input)
	if err != nil {
		return false, nil
	}
	return ok, nil
}

func UpperCaseOnly(input string) (bool, error) {
	pattern := "^[A-Z]+$"
	ok, err := regexp.MatchString(pattern, input)
	if err != nil {
		return false, nil
	}
	return ok, nil
}

func ComplexEngNum(input string) (bool, error) {
	pattern := "^[a-zA-Z0-9]+$"
	ok, err := regexp.MatchString(pattern, input)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func PersonNummber(input string) (bool, error) {
	pattern := "(^[0-9]{6})-([\\d]{7})$"
	ok, err := regexp.MatchString(pattern, input)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func Regex(pattern, input string) (bool, error) {
	ok, err := regexp.MatchString(pattern, input)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func main() {
	pattern := "(^[0-9]{6})-([12]{1}\\d{6})$"
	input := "910111-1909090"
	ok, _ := Regex(pattern, input)
	fmt.Println(ok)

	ok, _ = PersonNummber("909090-9090909")
	fmt.Println(ok)
	ok, _ = PersonNummber("90900-9090909")
	fmt.Println(ok)
	ok, _ = PersonNummber("90900-9090909")
	fmt.Println(ok)
}
