package main

import "fmt"

//_ "golearn/gohttp/sccservice"
//_ "golearn/gohttp/sccwork"

func add(base int) func(int) int {
	return func(i int) int {
		base += i
		return base
	}
}

func main() {
	tmp1 := make(map[string]interface{})
	tmp1["abc"] = 1
	tmp1["abd"] = 2
	tmp1["abe"] = 3
	tmp1["abf"] = 4
	tmp1["abg"] = 5
	tmp1["abh"] = 6
	tmp1["abi"] = 7
	for _, v := range tmp1 {
		fmt.Println(v)
	}
}
