package main

import "fmt"

func main() {
	c1 := make(chan interface{}, MaxInputDataLen)
	c2 := make(chan interface{}, MaxInputDataLen)
	c3 := make(chan interface{}, MaxInputDataLen)
	c4 := make(chan interface{}, MaxInputDataLen)
	c5 := make(chan interface{}, MaxInputDataLen)
	go Producer(c1, c2)
	go SingleHash(c2, c3)
	go MultiHash(c3, c4)
	go CombineResults(c4, c5)
	for val := range c5 {
		fmt.Println(val)
	}
}
