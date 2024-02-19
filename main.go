package main

func main() {
	flow := []job{
		producer,
		SingleHash,
	}
	ExecutePipeline(flow...)
}
