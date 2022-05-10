package example

import (
	"fmt"
	"github.com/dlclark/regexp2"
	"regexp"
	"strconv"
)

/*func Numeric1() {
	input := "9319"

	for _, i := range input {
		for j := 48; j <= 57; j++ {
			if int(i) == j {
				break
			}
			if j == 57 {
				fmt.Println("no numeric")
				break
			}
		}
	}
}*/

func Numeric1() {
	pattern := "^[a-zA-Z0-9]+$"
	input := "2rq112"
	ok, _ := regexp.MatchString(pattern, input)
	fmt.Println(ok) // true

	input = "2112"
	ok, _ = regexp.MatchString(pattern, input)
	fmt.Println(ok) // true

	input = "21WxZ12"
	ok, _ = regexp.MatchString(pattern, input)
	fmt.Println(ok) // true
}

func Numeric2() {
	input := "9319"
	digit := [10]int{}
	for _, i := range input {
		char := fmt.Sprintf("%c", i)

		for j := range digit {
			if char == strconv.Itoa(j) {
				break
			}
			if j == len(digit)-1 {
				fmt.Println("no numeric")
				break
			}
		}
	}
}

func Numeric3() {
	input := "93a19"
	digit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	for _, i := range input {
		char := fmt.Sprintf("%c", i)

		for j := range digit {
			if char == strconv.Itoa(j) {
				break
			}
			if j == len(digit)-1 {
				fmt.Println("no numeric")
				break
			}
		}
	}
}

func AlphaNumeric1() {
	input := "a2f@dz"

	for _, i := range input {
		for j := 48; j <= 57; j++ {
			if int(i) == j {
				break
			}
			if j == 57 {
				for k := 97; k <= 122; k++ {
					if int(i) == k {
						break
					}
					if k == 122 {
						fmt.Println("no alphanumeric")
						break
					}
				}
			}
		}
	}
}

func AlphaNumeric2() {
	input := "a2fTdz"

	for _, i := range input {
		for j := 48; j <= 57; j++ {
			if int(i) == j {
				break
			}
			if j == 57 {
				for k := 97; k <= 122; k++ {
					if int(i) == k {
						break
					}
					if k == 122 {
						for l := 65; l <= 90; l++ {
							if int(i) == l {
								break
							}
							if l == 90 {
								fmt.Println("no alphanumeric")
								break
							}
						}
					}
				}
			}
		}
	}
}

func RRN() {
	pattern := `([0-9]{6})-([12]{1}[\d]{6})`
	input := "881010-1090909"
	ok, _ := regexp.MatchString(pattern, input)
	fmt.Println(ok) // true

	input = "881010-5090909"
	ok, _ = regexp.MatchString(pattern, input)
	fmt.Println(ok) // false

	input = "8810101-509090"
	ok, _ = regexp.MatchString(pattern, input)
	fmt.Println(ok) // false
}

func PhoneNum() {
	pattern := `(010)-([\d]{4})-([\d]{4})`
	input := "010-2424-1111"
	ok, _ := regexp.MatchString(pattern, input)
	fmt.Println(ok) // true

	input = "011-2424-1111"
	ok, _ = regexp.MatchString(pattern, input)
	fmt.Println(ok) // false

	input = "010-424-1111"
	ok, _ = regexp.MatchString(pattern, input)
	fmt.Println(ok) // false
}

func CardNum() {
	pattern := `([\d]{4})-([\d]{4})-([\d]{4})-([\d]{4})`
	input := "0110-2424-1111-2411"
	ok, _ := regexp.MatchString(pattern, input)
	fmt.Println(ok) // true

	input = "011a-424-1111-2222"
	ok, _ = regexp.MatchString(pattern, input)
	fmt.Println(ok) // false

	input = "485-424-1111"
	ok, _ = regexp.MatchString(pattern, input)
	fmt.Println(ok) // false
}

func Email() {
	pattern := `^\w+([.-]?\w+)@\w+\.(\w{2,3})*.\w{2,3}`
	input := "abc@naver.co.kr"
	ok, _ := regexp.MatchString(pattern, input)
	fmt.Println(ok) // true

	input = "abc-@yahoo.com"
	ok, _ = regexp.MatchString(pattern, input)
	fmt.Println(ok) // false

	input = "ab.c@gmail.com"
	ok, _ = regexp.MatchString(pattern, input)
	fmt.Println(ok) // true

	input = "ab.c@gaggle.org"
	ok, _ = regexp.MatchString(pattern, input)
	fmt.Println(ok) // true

	input = "abc@naver.co.."
	ok, _ = regexp.MatchString(pattern, input)
	fmt.Println(ok) // false
}

func Password10() {
	pattern := `^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!@#$%^&*\(\)\-_=+])([a-zA-Z0-9!@#$%^&*\(\)\-_=+]+){10,}$`
	//pattern := `^(?=.{0,}[a-z])(?=.{0,}[A-Z])(?=.{0,}[0-9])(?=.{0,}[!@#$%^&*\(\)\-_=+])([a-zA-Z0-9!@#$%^&*\(\)\-_=+]+){10,}$`
	re := regexp2.MustCompile(pattern, regexp2.RE2)

	input := "Awdfgesd#fDS"
	ok, _ := re.MatchString(input)
	fmt.Println(ok) // false

	input = "#Wfg3q"
	ok, _ = re.MatchString(input)
	fmt.Println(ok) // false

	input = "weE2qacvsq+wQsddsdW@!"
	ok, _ = re.MatchString(input)
	fmt.Println(ok) // true

	input = "dfwgD21qasd"
	ok, _ = re.MatchString(input)
	fmt.Println(ok) // false
}

func Profanity() {
	re := regexp.MustCompile(`ㅅㅂ|ㅄ`)
	msg := "ㅅㅂ아 ~~ ㅄ아 ~~"
	input := re.ReplaceAllString(msg, "아잉")
	fmt.Println(input) // 아잉아 ~~ 아잉아 ~~
}
