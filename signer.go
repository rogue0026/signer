package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// сюда писать код

var jobsFlow = []job{
	func(in, out chan interface{}) {
		inputData := []int{0, 1, 1, 2, 3, 5, 8}
		for _, elem := range inputData {
			out <- elem
		}
		//close(out)
	},
	SingleHash,
	MultiHash,
	CombineResults,
	func(in, out chan interface{}) {
		val := <-in
		if res, ok := val.(string); ok {
			fmt.Println(res)
		}
	},
}

var ExecutePipeline = func(jobFunctions ...job) {
	firstGwg := sync.WaitGroup{}
	globWg := sync.WaitGroup{}
	channels := make([]chan interface{}, len(jobFunctions))
	for i := 0; i < len(channels); i++ {
		channels[i] = make(chan interface{}, MaxInputDataLen)
	}

	for i := 0; i < len(jobFunctions); i++ {
		firstGoroutine := i == 0

		if firstGoroutine {
			firstGwg.Add(1)
			go func(wg *sync.WaitGroup, th int) {
				defer wg.Done()
				jobFunctions[th](nil, channels[th])
			}(&firstGwg, i)
		} else {
			// wait for results from goroutines
			globWg.Add(1)
			go func(wg *sync.WaitGroup, th int) {
				defer globWg.Done()
				jobFunctions[th](channels[th-1], channels[th])
			}(&globWg, i)
		}
	}

	go func(wg *sync.WaitGroup) {
		wg.Wait()
		close(channels[0])
	}(&firstGwg)

	globWg.Wait()
}

// SingleHash считает значение crc32(data)+"~"+crc32(md5(data)) ( конкатенация двух строк через ~),
// где data - то что пришло на вход (по сути - числа из первой функции)
var SingleHash = func(in, out chan interface{}) {
	wg := sync.WaitGroup{}
	for val := range in {
		if n, ok := val.(int); ok {
			data := strconv.Itoa(n)
			hashData := DataSignerMd5(data)

			wg.Add(1)
			go func(normalData string, hash string) {
				defer wg.Done()
				hashCRCCh := make(chan string, 1)
				dataCRCCh := make(chan string, 1)
				go func(hashCRCCh chan string) {
					hashCRC := DataSignerCrc32(hash)
					hashCRCCh <- hashCRC
				}(hashCRCCh)
				go func(dataCRCCh chan string) {
					dataCRC := DataSignerCrc32(normalData)
					dataCRCCh <- dataCRC
				}(dataCRCCh)
				dCRC := <-dataCRCCh
				hCRC := <-hashCRCCh
				out <- dCRC + "~" + hCRC
			}(data, hashData)
		}
	}
	wg.Wait()
	close(out)
	// дождаться пока все отправят данные
}

// MultiHash считает значение crc32(th+data) (конкатенация цифры, приведённой к строке и строки),
// где th=0..5 ( т.е. 6 хешей на каждое входящее значение), потом берёт конкатенацию результатов в порядке расчета
// (0..5), где data - то что пришло на вход (и ушло на выход из SingleHash)
var MultiHash = func(in, out chan interface{}) {
	wg := sync.WaitGroup{}
	for val := range in {
		results := make([]chan string, 6)
		for i := 0; i < len(results); i++ {
			results[i] = make(chan string, 1)
		}
		if singHashVal, ok := val.(string); ok {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for i := 0; i < 6; i++ {
					go func(th int) {
						results[th] <- DataSignerCrc32(strconv.Itoa(th) + singHashVal)
					}(i)
				}
				multiHashParts := make([]string, 6)
				for i := 0; i < len(results); i++ {
					multiHashParts[i] = <-results[i]
				}
				out <- strings.Join(multiHashParts, "")
			}()
		}
	}
	// wait here while all goroutines send data
	wg.Wait()
	close(out)
}

// CombineResults получает все результаты, сортирует (https://golang.org/pkg/sort/),
// объединяет отсортированный результат через _ (символ подчеркивания) в одну строку
var CombineResults = func(in, out chan interface{}) {
	results := make([]string, 0, MaxInputDataLen)
	for val := range in {
		if multiHashVal, ok := val.(string); ok {
			results = append(results, multiHashVal)
		}
	}
	sort.Slice(results, func(i int, j int) bool { return results[i] < results[j] })
	out <- strings.Join(results, "_")
	close(out)
}
