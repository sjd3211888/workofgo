package main

import "fmt"

type aa struct {
	xxx int
}
type bb struct {
	abc map[int]*aa
	m   int
}
type cc struct {
	bcd map[int]*aa
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	var ch chan int = make(chan int, 10)
	close(ch)
	ch <- 1
}
