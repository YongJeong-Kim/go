package ch

import "fmt"

func Ex1() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	defer close(ch1)
	defer close(ch2)

	ch1 = sum1(ch1)
	ch2 = sub1(ch2, <-ch1)

	fmt.Println(<-ch2)
}

func sum1(ch chan int) chan int {
	var s int

	go func() {
		for i := 1; i <= 1000; i++ {
			s += i
		}
		ch <- s
	}()

	return ch
}

func sub1(sch chan int, num int) chan int {
	if num == 500500 {
		go func() {
			for i := 1000; i >= 1; i-- {
				num -= i
			}
			sch <- num
		}()
	} else {
		go func() {
			sch <- 111
		}()
	}

	return sch
}