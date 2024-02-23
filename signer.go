package main

import (
	"strconv"
	"strings"
	"sync"
)

// сюда писать код

// crc(data) + crc(hash(data))
var SingleHash = func(in, out chan interface{}) {
	wg := sync.WaitGroup{}
	//m := sync.Mutex{}
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

var CombineResults = func(in, out chan interface{}) {}
