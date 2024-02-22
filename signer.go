package main

import (
	"sort"
	"strconv"
	"strings"
)

// сюда писать код
var jobsFlow = []job{
	Producer,
	SingleHash,
	MultiHash,
	CombineResults,
}

func ExecutePipeline(jobs ...job) {
	//for i := 0; i < len(jobs); i++ {
	//	jobs[i]()
	//}
}

func Producer(in, out chan interface{}) {
	inputData := []int{0, 1, 1, 2, 3, 5, 8}
	for _, elem := range inputData {
		out <- elem
	}
	close(out)
}

func SingleHash(in, out chan interface{}) {
	for val := range in {
		if num, ok := val.(int); ok {
			hash := DataSignerMd5(strconv.Itoa(num))
			hashCRC := DataSignerCrc32(hash)
			dataCRC := DataSignerCrc32(strconv.Itoa(num))
			out <- dataCRC + "~" + hashCRC
		}
	}
	close(out)
}

func MultiHash(in, out chan interface{}) {
	for val := range in {
		var multiHashVal string
		if singHash, ok := val.(string); ok {
			for i := 0; i < 6; i++ {
				multiHashVal += DataSignerCrc32(strconv.Itoa(i) + singHash)
			}
			out <- multiHashVal
		}
	}
	close(out)
}

func CombineResults(in, out chan interface{}) {
	multiHashResults := make([]string, 0, MaxInputDataLen)
	for val := range in {
		if multiHashRes, ok := val.(string); ok {
			multiHashResults = append(multiHashResults, multiHashRes)
		}
	}
	sort.Slice(multiHashResults, func(i int, j int) bool { return multiHashResults[i] < multiHashResults[j] })
	out <- strings.Join(multiHashResults, "_")
	close(out)
}
