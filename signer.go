package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

func ExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}
	in := make(chan interface{})
	out := make(chan interface{})
	for _, j := range jobs {
		wg.Add(1)
		go func(j job) {
			defer wg.Done()
			j(in, out)
		}(j)
		in = out
		out = make(chan interface{})
	}
	wg.Wait()
}

func producer(in chan interface{}, out chan interface{}) {
	inputData := []int{0, 1, 1, 2, 3, 5}
	for _, val := range inputData {
		out <- val
		time.Sleep(time.Millisecond * 10)
	}
	close(out)
}

// crc32(data)+"~"+crc32(md5(data))
func SingleHash(in chan interface{}, out chan interface{}) {
	for rawData := range in {
		num, ok := rawData.(int)
		if ok {
			data := strconv.Itoa(num)
			c := make(chan string, 1) // этот канал связывает две горутины: первая вычисляет md5, вторая вычисляет crc32 на основе md5

			go func() {
				res := DataSignerCrc32(data)
				fmt.Printf("crc32: %v\n", res)
				out <- res
			}()

			go func(output chan string) {
				output <- DataSignerMd5(data)
				close(output)
			}(c)
			go func(input chan string) {
				res := DataSignerCrc32(<-input)
				fmt.Printf("HashCrc32: %v\n", res)
				out <- res
			}(c)
		}
	}
	close(out)
}

func MultiHash(in chan interface{}, out chan interface{})      {}
func CombineResults(in chan interface{}, out chan interface{}) {}
