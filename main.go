package main

import "time"

func main() {
	pipeline := []job{gen, SingleHash, MultiHash, CombineResults}
	ExecutePipeline(pipeline...)
	time.Sleep(time.Second * 20)
}
