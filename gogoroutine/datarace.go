package main

import (
	"fmt"
	"time"
)

type account struct {
	money int
}

func datarace() {
	acc := &account{
		money: 10000,
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
		}
	}

	for i := 0; i < 50; i++ {
		go func() {
			work(acc)
		}()
	}

	for i := 0; i < 50; i++ {
		time.Sleep(time.Millisecond * 100)
		fmt.Println(acc.money)
	}
}