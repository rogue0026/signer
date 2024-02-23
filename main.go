package main

func main() {
	//c1 := make(chan interface{}, 10)
	//c2 := make(chan interface{}, 10)
	//c3 := make(chan interface{}, 10)
	//c4 := make(chan interface{}, 10)
	//inpData := []int{0, 1, 1, 2, 3, 5, 8}
	//go func() {
	//	for _, elem := range inpData {
	//		c1 <- elem
	//	}
	//	close(c1)
	//}()
	//SingleHash(c1, c2)
	//MultiHash(c2, c3)
	//CombineResults(c3, c4)
	//for res := range c4 {
	//	fmt.Println(res)
	//}
	ExecutePipeline(jobsFlow...)
}
