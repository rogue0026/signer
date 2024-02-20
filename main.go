package main

import (
	"fmt"
	"os"
)

var jobsFlow = []job{
	func(in chan interface{}, out chan interface{}) {
		inputData := []int{0, 1, 1, 2, 3, 5}
		for _, elem := range inputData {
			out <- elem
		}
		close(out)
	},
	SingleHash,
	MultiHash,
	CombineResults,
	func(in chan interface{}, out chan interface{}) {
		for rawVal := range in {
			if res, ok := rawVal.(string); ok {
				fmt.Fprintln(os.Stdout, res)
			}
		}
		out <- struct{}{}
		close(out)
	},
}

func main() {
	ExecutePipeline(jobsFlow...)
	fmt.Scanln()
	// ch1 := make(chan interface{})
	// ch2 := make(chan interface{})
	// ch3 := make(chan interface{})
	// go numGenerator(ch1)
	// go SingleHash(ch1, ch2)
	// go MultiHash(ch2, ch3)

	// for elem := range ch3 {
	// 	fmt.Println(elem)
	// }
}

// func numGenerator(out chan interface{}) {
// 	inputData := []int{0, 1, 1, 2, 3, 5}
// 	for _, elem := range inputData {
// 		out <- elem
// 	}
// 	close(out)
// }
