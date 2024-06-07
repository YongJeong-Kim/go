package service

/*func TestAd(t *testing.T) {
	as := &asdf{
		a:  make(chan string),
		wg: sync.WaitGroup{},
	}
	//as.wg.Add(1)

	go as.aa()
	go as.cc()
	as.bb()
	//as.wg.Wait()
}

type asdf struct {
	a  chan string
	wg sync.WaitGroup
}

func (a *asdf) aa() {
	//a.wg.Add(1)
	//defer a.wg.Done()
	a.a <- "hello"
	a.a <- "hello!!"
}

func (a *asdf) cc() {
	a.wg.Add(1)
	defer a.wg.Done()
	a.a <- "hello22"
	a.a <- "hello22!!"
}

func (a *asdf) bb() {
	for {
		select {
		case se := <-a.a:
			fmt.Println(se)
		}
	}
}
*/
