package main

import "fmt"

func main() {
	jobs := []job{
		func(in, out chan interface{}) {
			inputData := []int{0, 1, 1, 2, 3, 5}
			for _, data := range inputData {
				out <- data
			}
			close(out)
		},
		SingleHash,
		MultiHash,
		func(in, out chan interface{}) {
			for data := range in {
				val, ok := data.(string)
				if ok {
					fmt.Println(val)
				}
			}
		},
	}
	ExecutePipeline(jobs...)
	fmt.Scanln()
	// ch1 := make(chan interface{})
	// ch2 := make(chan interface{})
	// go numGenerator(ch1)
	// go SingleHash(ch1, ch2)
	// for elem := range ch2 {
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
