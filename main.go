package main

import (
	"fmt"
)

func main() {
	jobsFlow := []job{
		SingleHash,
		MultiHash,
		CombineResults,
	}
	ExecutePipeline(jobsFlow...)
	fmt.Scanln()
}

func receiver(in chan interface{}, out chan interface{}) {
	for val := range in {
		fmt.Println(val)
	}
}

func generator(in chan interface{}, out chan interface{}) {
	inputData := []int{0, 1, 1, 2, 3, 5, 8}
	for _, val := range inputData {
		out <- val
	}
	close(out)
}
