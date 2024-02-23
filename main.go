package main

import "fmt"

func main() {
	c1 := make(chan interface{}, 10)
	c2 := make(chan interface{}, 10)
	inpData := []int{0, 1, 1, 2, 3, 5, 8}
	go func() {
		for _, elem := range inpData {
			c1 <- elem
		}
		close(c1)
	}()
	SingleHash(c1, c2)

	for res := range c2 {
		fmt.Println(res)
	}
}
