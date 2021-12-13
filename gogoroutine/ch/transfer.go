package ch

import (
	"fmt"
	"time"
)

type account struct {
	money int
	ch chan int
}

func Transfer() {
	 acc := &account{
		money: 10000,
		ch: make(chan int),
	 }

	deposit := func(a *account) {
		a.money += 1000
	}

	withdraw := func(a *account) {
		a.money -= 1000
	}

	work := func(a *account) {
		for {
			deposit(a)
			withdraw(a)
			a.ch <- a.money
		}
	}

	for i := 0; i < 20; i++ {
		go func() {
			work(acc)
		}()
	}

	for i := 0; i < 50; i++ {
		acc.money = <-acc.ch
		time.Sleep(time.Millisecond * 100)
		fmt.Println(acc.money)
	}
}