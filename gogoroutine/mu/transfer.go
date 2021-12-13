package mu

import (
	"fmt"
	"sync"
	"time"
)

type account struct {
	money int
	mutex *sync.Mutex
}

func Transfer() {
	acc := &account{
		money: 10000,
		mutex: &sync.Mutex{},
	}

	deposit := func(a *account) {
		a.money += 1000
	}

	withdraw := func(a *account) {
		a.money -= 1000
	}

	work := func(a *account) {
		for {
			a.mutex.Lock()
			deposit(a)
			withdraw(a)
			a.mutex.Unlock()
		}
	}

	for i := 0; i < 20; i++ {
		go func() {
			work(acc)
		}()
	}

	for i := 0; i < 50; i++ {
		time.Sleep(time.Millisecond * 100)
		acc.mutex.Lock()
		fmt.Println(acc.money)
		acc.mutex.Unlock()
	}
}