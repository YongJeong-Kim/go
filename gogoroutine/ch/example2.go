package ch

import "fmt"

func Ex2() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	defer close(ch1)
	defer close(ch2)

	go sum2(ch1)
	go sub2(ch2, ch1)

	fmt.Println(<-ch2)
}

func sum2(ch chan int) {
	var s int

	for i := 1; i <= 1000; i++ {
		s += i
	}
	ch <- s
}

func sub2(sch chan<- int, rch <-chan int) {
	num := <-rch
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
}