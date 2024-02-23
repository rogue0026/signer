package main

import "fmt"

func main() {
	c1 := make(chan interface{}, 10)
	c2 := make(chan interface{}, 10)
	c3 := make(chan interface{}, 10)
	inpData := []int{0, 1, 1, 2, 3, 5, 8}
	go func() {
		for _, elem := range inpData {
			c1 <- elem
		}
		close(c1)
	}()
	SingleHash(c1, c2)
	MultiHash(c2, c3)
	for res := range c3 {
		fmt.Println(res)
	}
}
